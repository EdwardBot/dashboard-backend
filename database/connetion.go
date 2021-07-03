package database

import (
	"context"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"log"
	"os"
)

var (
	Conn *pgx.Conn
)

func Connect() {
	conn, err := pgx.ConnectConfig(context.Background(), &pgx.ConnConfig{
		Config: pgconn.Config{
			Host:           "45.135.56.198",
			Port:           5432,
			Database:       "edward",
			User:           "admin",
			Password:       os.Getenv("DB_PASS"),
			ConnectTimeout: 10000,
			AfterConnect:   nil,
		},
	})
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}

	Conn = conn

	defer func(conn *pgx.Conn, ctx context.Context) {
		err := conn.Close(ctx)
		if err != nil {
			log.Println(err)
		}
	}(conn, context.Background())
}
