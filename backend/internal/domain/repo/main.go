package repo

import (
	"context"

	"github.com/densmart/smart-stream/internal/adapters/db/postgres"
	"github.com/densmart/smart-stream/pkg/logger"
)

type OltpRepo interface {
	Close()
	MigrationUp() error
	MigrationDown() error
	Clients
}

func NewOltpRepo(ctx context.Context, DBType string) (OltpRepo, error) {
	var database OltpRepo
	switch DBType {
	default:
		pg, err := postgres.NewPostgresDB(ctx)
		if err != nil {
			return nil, err
		}
		database = pg
	}
	// run migrations
	if err := database.MigrationUp(); err != nil {
		logger.Fatalf("error migrating DB: %s", err.Error())
	}

	return database, nil
}
