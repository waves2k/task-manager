package database

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	TodoTable  = "todos_user"
	UsersTable = "users"
)

func Connect(connectionString string) (*pgxpool.Pool, error) {
	ctx := context.Background()

	config, err := pgxpool.ParseConfig(connectionString)
	if err != nil {
		log.Printf("Unable to parse CONNECTION_STRING: %V", err)
		return nil, err
	}
	pool, err := pgxpool.NewWithConfig(ctx, config)

	if err != nil {
		log.Printf("Unable to create connection pool %v", err)
		return nil, err
	}
	if err = pool.Ping(ctx); err != nil {
		log.Printf("Unable to ping database: %v", err)
		pool.Close()
		return nil, err
	}

	log.Printf("Succesfully connected to PostgreSQL database")
	return pool, nil
}
