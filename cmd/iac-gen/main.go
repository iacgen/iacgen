package main

import (
	"context"
	"os"
	"os/signal"
	"time"

	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"

	"github.com/cafi-dev/iac-gen/pkg/httpserver"
	"github.com/cafi-dev/iac-gen/pkg/logging"
)

func main() {
	logger := logging.GetLogger()
	logger.Info("starting server")

	// create context
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer cancel()

	// error groups for running multiple functionalites concurrently
	errGrp, ctx := errgroup.WithContext(ctx)

	// start http server
	s := httpserver.NewHTTPServer("8000")
	errGrp.Go(func() error {
		return s.Start()
	})

	errGrp.Go(func() error {
		<-ctx.Done()
		ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
		defer cancel()
		if err := s.Shutdown(ctx); err != nil {
			logger.Warn("failed to shutdown http server", zap.Error(err))
		}
		return ctx.Err()
	})

	err := errGrp.Wait()
	logger.Error("stopping server", zap.Error(err))

	if err := s.Close(); err != nil {
		logger.Error("error closing http server")
	}
	logger.Info("service terminated!")
}
