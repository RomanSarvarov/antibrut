package main

import (
	"context"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"

	"antibrut/config"
	"antibrut/leakybucket"
	"antibrut/proto/antibrut/v1"
	"antibrut/sqlite"
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

		if err = db.Migrate(cfg.SQLite.MigrationsPath); err != nil {
			return err
		}
	}

	// rate limiter
	limiter := leakybucket.New(db)
	_ = limiter

	// grpc server
	lis, err := net.Listen("tcp", cfg.GRPC.Address)
	if err != nil {
		return errors.Wrap(err, "start grpc server error")
	}

	srv := grpc.NewServer()

	antibrut.RegisterAntiBrutServiceServer(srv, grpcapi.New(model))

	closer.Add(func() error {
		log.
			Debug().
			Msgf("terminating GRPC server")

		srv.GracefulStop()

		return nil
	})

	errgrp.Go(func() error {
		log.
			Debug().
			Msgf("starting GRPC server on: `%s`", cfg.GRPC.Address)

		err := srv.Serve(lis)
		if err != nil && !errors.Is(err, grpc.ErrServerStopped) {
			return err
		}
		return nil
	})

	<-ctx.Done()

	cmd.Println("Stopping...")

	if err = errGrp.Wait(); err != nil {
		return err
	}

	return nil
}
