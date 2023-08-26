package main

import (
	"context"
	"dynamic-user-segmentation/config"
	"dynamic-user-segmentation/internal/ports/httpgin"
	segmentRepo "dynamic-user-segmentation/internal/repository/segment"
	usersSegmentRepo "dynamic-user-segmentation/internal/repository/user_segment"
	segmentService "dynamic-user-segmentation/internal/service/segment"
	userSegmentService "dynamic-user-segmentation/internal/service/user_segment"
	"dynamic-user-segmentation/pkg/logging"
	"dynamic-user-segmentation/pkg/postgres"
	"errors"
	"fmt"
	"golang.org/x/sync/errgroup"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		panic(fmt.Errorf("unable to load config: %w", err))
	}

	log, err := logging.New(cfg.Logger)
	if err != nil {
		panic(fmt.Errorf("unable to configure logger: %w", err))
	}

	log.Info("connecting to postgres...")
	pgPool, err := postgres.New(cfg.DSN, 5, log)
	if err != nil {
		log.Fatal(fmt.Errorf("unable to connect to postgres: %w", err))
	}
	defer pgPool.Close()
	log.Info("successfully connected")

	log.Info("configure server...")
	server := httpgin.NewServer(
		":"+cfg.Server.Port,
		segmentService.New(segmentRepo.New(pgPool)),
		userSegmentService.New(usersSegmentRepo.New(pgPool)),
		log)
	log.Info("successfully configured")

	g, ctx := errgroup.WithContext(context.Background())
	gracefulShutdown(ctx, g)

	g.Go(func() error {
		log.Infof("starting http server on port: %s\n", server.Addr)
		defer log.Infof("closing http server on port: %s\n", server.Addr)

		errCh := make(chan error)

		defer func() {
			shCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()

			if err := server.Shutdown(shCtx); err != nil {
				log.Infof("can't close http server listening on %s: %s", server.Addr, err.Error())
			}

			close(errCh)
		}()

		go func() {
			if err = server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
				errCh <- err
			}
		}()
		select {
		case <-ctx.Done():
			return ctx.Err()
		case err = <-errCh:
			return fmt.Errorf("http server can't listen and serve requests: %w", err)
		}
	})

	if err = g.Wait(); err != nil {
		log.Infof("gracefully shutting down the server: %s\n", err.Error())
	}
}

func gracefulShutdown(ctx context.Context, g *errgroup.Group) {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
	g.Go(func() error {
		select {
		case s := <-signals:
			return fmt.Errorf("captured signal %s\n", s)
		case <-ctx.Done():
			return nil
		}
	})
}
