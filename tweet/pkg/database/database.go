package database

import (
	"context"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/jackgris/twitter-backend/tweet/pkg/logger"
)

func ConnectDB(ctx context.Context, log *logger.Logger) *pgx.Conn {
	conn, err := pgx.Connect(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Error(ctx, "Unable to connect to database", "Can't open DB connection", err)
		os.Exit(1)
	}

	err = conn.Ping(ctx)
	if err != nil {
		log.Error(ctx, "Connection database", "ERROR", err)
		os.Exit(1)
	}

	return conn
}
