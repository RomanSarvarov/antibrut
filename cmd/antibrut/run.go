package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"

	"github.com/romsar/antibrut"
	"github.com/romsar/antibrut/config"
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

	// rate limiter
	rateLimiter := leakybucket.New(db)

	// service
	service := antibrut.NewService(db, rateLimiter)

	// grpc server
	server := localgrpc.NewServer(service)

	errGrp.Go(func() error {
		cmd.Printf("Starting server on `%s`\n", cfg.GRPC.Address)

		if err := server.Start(cfg.GRPC.Address); err != nil {
			return err
		}
		return nil
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
