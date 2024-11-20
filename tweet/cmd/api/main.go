package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/jackgris/twitter-backend/tweet/internal/handler"
	"github.com/jackgris/twitter-backend/tweet/internal/store/tweetdb"
	"github.com/jackgris/twitter-backend/tweet/pkg/database"
	"github.com/jackgris/twitter-backend/tweet/pkg/logger"
)

func main() {
	ctx := context.Background()
	log := logger.New(os.Stdout)

	err := run(ctx, log)
	if err != nil {
		log.Error(ctx, fmt.Sprintf("Error server shutdown: %s\n", err))
	}
}

func run(ctx context.Context, log *logger.Logger) error {

	db := database.ConnectDB(ctx, log)
	defer db.Close(ctx)

	store := tweetdb.NewStore(db)
	mux := handler.NewHandler(store, log)

	// portEnv := os.Getenv("PORT")
	// port, err := strconv.Atoi(portEnv)
	// if err != nil {
	// 	log.Error(ctx, "Environment variable PORT converting to integer")
	// 	os.Exit(1)
	// }
	port := "8083"

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	serverErrors := make(chan error, 1)
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Info(ctx, "startup", "GOMAXPROCS", runtime.GOMAXPROCS(0), "Server started at port", port)

		serverErrors <- srv.ListenAndServe()
	}()

	// -------------------------------------------------------------------------
	// Shutdown

	select {
	case err := <-serverErrors:
		return fmt.Errorf("server error: %w", err)

	case sig := <-shutdown:
		log.Info(ctx, "shutdown", "status", "shutdown started", "signal", sig)
		defer log.Info(ctx, "shutdown", "status", "shutdown complete", "signal", sig)

		ctx, cancel := context.WithTimeout(ctx, time.Microsecond*500)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			return fmt.Errorf("could not stop server gracefully: %w", err)
		}
	}

	return nil
}
