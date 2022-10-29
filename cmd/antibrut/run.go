package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"

	"github.com/romsar/antibrut"
	"github.com/romsar/antibrut/config"
	"github.com/romsar/antibrut/inmem"
	"github.com/romsar/antibrut/leakybucket"
	"github.com/romsar/antibrut/sqlite"

	localgrpc "github.com/romsar/antibrut/grpc"
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Запускает HTTP сервер и начинает ожидать входящие вызовы.",
	RunE:  serve,
}

func init() {
	rootCmd.AddCommand(runCmd)

	runCmd.Flags().Bool("migrations", true, "Run database migrations before start.")
}

func serve(cmd *cobra.Command, args []string) error {
	// flags
	runMs, err := cmd.Flags().GetBool("migrations")
	if err != nil {
		return err
	}

	// config
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	// graceful shutdown
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	errGrp, ctx := errgroup.WithContext(ctx)

	// db
	cmd.Println("Connecting to database...")

	db, err := sqlite.New(cfg.SQLite.DSN)
	if err != nil {
		return err
	}
	defer db.Close()

	if runMs {
		cmd.Println("Running database migrations...")

		if err := db.Migrate(); err != nil {
			return err
		}
	}

	var lbRepo leakybucket.Repository

	// inmem
	switch cfg.RateLimiterStorageDriver {
	case "sqlite":
		lbRepo = db
	case "inmem":
		lbRepo = buildInMemoryLBRepository(db)
	default:
		return errors.New("unknown ratelimiter driver")
	}

	// rate limiter
	cmd.Printf(
		"Using Leaky Bucket rate limiter with `%s` storage driver.\n",
		cfg.RateLimiterStorageDriver,
	)

	rateLimiter := leakybucket.New(lbRepo)

	// service
	service := antibrut.NewService(
		db,
		rateLimiter,
		antibrut.WithPruneDuration(cfg.PruneDuration),
	)

	// grpc server
	server := localgrpc.NewServer(service)

	errGrp.Go(func() error {
		defer cancel()

		cmd.Printf("Starting server on `%s`\n", cfg.GRPC.Address)

		if err := server.Start(cfg.GRPC.Address); err != nil {
			return err
		}
		return nil
	})

	// antibrut worker
	errGrp.Go(func() error {
		defer cancel()

		return service.Work(ctx)
	})

	<-ctx.Done()

	cmd.Println("Stopping...")

	if err := server.Close(); err != nil {
		cmd.PrintErrln(err)
	}

	if err := errGrp.Wait(); err != nil {
		return err
	}

	return nil
}

// InMemLBRepository in-memory репозиторий для Leaky Bucket алгоритма.
// Все операции, кроме поиска antibrut.Limitation происходят в памяти.
type InMemLBRepository struct {
	db *sqlite.Repository
	*inmem.Repository
}

// FindLimitation находит antibrut.Limitation.
// Если совпадений нет, вернет antibrut.ErrNotFound.
func (r *InMemLBRepository) FindLimitation(ctx context.Context, c antibrut.LimitationCode) (*antibrut.Limitation, error) {
	return r.db.FindLimitation(ctx, c)
}

// buildInMemoryLBRepository создает InMemLBRepository.
func buildInMemoryLBRepository(db *sqlite.Repository) leakybucket.Repository {
	return &InMemLBRepository{
		Repository: inmem.New(),
		db:         db,
	}
}
