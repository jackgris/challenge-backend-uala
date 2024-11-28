package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"strconv"
	"syscall"
	"time"

	"github.com/jackgris/twitter-backend/timeline/internal/handler"
	"github.com/jackgris/twitter-backend/timeline/internal/store/timelinedb"
	"github.com/jackgris/twitter-backend/timeline/pkg/database"
	"github.com/jackgris/twitter-backend/timeline/pkg/logger"
	"github.com/jackgris/twitter-backend/timeline/pkg/msgbroker"
)

func main() {
	ctx := context.Background()
	log := logger.New(os.Stdout)

	err := run(ctx, log)
	if err != nil {
		log.Error(ctx, "timeline service", fmt.Sprintf("Error server shutdown: %s\n", err))
	}
}

func run(ctx context.Context, log *logger.Logger) error {

	db := database.ConnectDB(ctx, log)
	defer db.Close(ctx)

	store := timelinedb.NewStore(db)
	msgBrokerPath := os.Getenv("NATS_URL")
	if msgBrokerPath == "" {
		log.Error(ctx, "timeline service", "Environment variable NATS_URL is empty")
		os.Exit(1)
	}

	msgbroker := msgbroker.NewMsgBroker(msgBrokerPath, log)
	mux := handler.NewHandler(store, msgbroker, log)

	portEnv := os.Getenv("PORT")
	port, err := strconv.Atoi(portEnv)
	if err != nil {
		log.Error(ctx, "timeline service", "Environment variable PORT converting to integer", err)
		os.Exit(1)
	}

	srv := &http.Server{
		Addr:    ":" + strconv.Itoa(port),
		Handler: mux,
	}

	serverErrors := make(chan error, 1)
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Info(ctx, "timeline service startup", "GOMAXPROCS", runtime.GOMAXPROCS(0), "Server started at port", port)

		serverErrors <- srv.ListenAndServe()
	}()

	// -------------------------------------------------------------------------
	// Shutdown

	select {
	case err := <-serverErrors:
		log.Error(ctx, "timeline service", "status", "server error", err)
		return fmt.Errorf("server error: %w", err)

	case sig := <-shutdown:
		log.Info(ctx, "timeline service shutdown", "status", "shutdown started", "signal", sig)
		defer log.Info(ctx, "timeline service shutdown", "status", "shutdown complete", "signal", sig)

		ctx, cancel := context.WithTimeout(ctx, time.Microsecond*500)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			log.Error(ctx, "timeline service shutdown", "status", "could not stop server gracefully", err)
			return fmt.Errorf("could not stop server gracefully: %w", err)
		}
	}

	return nil
}
