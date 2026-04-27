package postgres

import (
	"context"

	"github.com/densmart/smart-stream/pkg/logger"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/jackc/tern/migrate"
	"github.com/spf13/viper"
)

const (
	migrationsDirectory = "./backend/migrations"
)

type PgDB struct {
	ctx  context.Context
	pool *pgxpool.Pool
}

func NewPostgresDB(ctx context.Context) (*PgDB, error) {
	connString := viper.GetString("db.postgres.dsn")
	logger.Debugf("connecting to DB: %s", connString)
	conf, err := pgxpool.ParseConfig(connString) // Using environment variables instead of a connection string.
	if err != nil {
		logger.Errorf("%s", err.Error())
		return nil, err
	}

	conf.ConnConfig.LogLevel = pgx.LogLevelWarn
	conf.MaxConns = 50
	conf.ConnConfig.PreferSimpleProtocol = true

	pool, err := pgxpool.ConnectConfig(ctx, conf)
	if err != nil {
		logger.Errorf("%s", err.Error())
		return nil, err
	}

	if err = getConnection(ctx, pool); err != nil {
		logger.Errorf("%s", err.Error())
		return nil, err
	}

	return &PgDB{
		ctx:  ctx,
		pool: pool,
	}, nil
}

func (db *PgDB) Close() {
	db.pool.Close()
}

func (db *PgDB) MigrationUp() error {
	conn, err := db.pool.Acquire(db.ctx)
	defer conn.Release()
	if err != nil {
		return err
	}
	migrator, err := migrate.NewMigrator(context.Background(), conn.Conn(), "schema_version")
	if err != nil {
		return err
	}
	if err = migrator.LoadMigrations(migrationsDirectory); err != nil {
		return err
	}
	if err = migrator.Migrate(context.Background()); err != nil {
		return err
	}
	return nil
}

func (db *PgDB) MigrationDown() error {
	return nil
}

// get connection from pool and release
func getConnection(ctx context.Context, pool *pgxpool.Pool) error {
	conn, err := pool.Acquire(ctx)
	defer conn.Release()
	if err != nil {
		logger.Errorf("%s", err.Error())
		return err
	}
	if err = conn.Ping(ctx); err != nil {
		logger.Errorf("%s", err.Error())
		return err
	}
	return nil
}
