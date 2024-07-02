package main

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/quolpr/wisdom-pow/internal/service/pow"
	"github.com/quolpr/wisdom-pow/internal/service/quotes"
	"github.com/quolpr/wisdom-pow/internal/service/quotes/repo/jsonquote"
	"github.com/quolpr/wisdom-pow/internal/tcpserver"
	"github.com/quolpr/wisdom-pow/internal/tcptransport"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		<-c
		cancel()

		slog.Info("Shutting down...")
	}()

	powService := pow.NewService()
	repo, err := jsonquote.NewRepo()
	if err != nil {
		slog.ErrorContext(ctx, "Error creating quote service:", "error", err)

		os.Exit(1) //nolint:gocritic
	}
	quoteService := quotes.NewService(repo)

	err = tcpserver.StartServer(ctx, 8080, tcptransport.NewHandler(powService, quoteService, 20))

	if err != nil && !errors.Is(err, context.Canceled) {
		slog.ErrorContext(ctx, "Error starting server:", "error", err)
		os.Exit(1)
	}
}
