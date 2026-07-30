package postgres

import (
	"context"
	"fmt"
	"strings"

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
	// Try to get full DSN first, otherwise build it from individual parameters
	connString := viper.GetString("db.postgres.dsn")
	if connString == "" {
		// Build DSN from individual parameters
		host := viper.GetString("db.postgres.host")
		port := viper.GetInt("db.postgres.port")
		database := viper.GetString("db.postgres.database")
		user := viper.GetString("db.postgres.user")
		password := viper.GetString("db.postgres.password")
		sslmode := viper.GetString("db.postgres.sslmode")
		connString = buildDSN(host, port, database, user, password, sslmode)
	}
	logger.Debugf("connecting to DB: %s", maskPassword(connString))
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

// buildDSN constructs PostgreSQL connection string from individual parameters
func buildDSN(host string, port int, database, user, password, sslmode string) string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		user, password, host, port, database, sslmode)
}

// maskPassword masks password in DSN for logging
func maskPassword(dsn string) string {
	// Find password part: between "://" and "@"
	if idx := strings.Index(dsn, "://"); idx != -1 {
		if atIdx := strings.Index(dsn[idx+3:], "@"); atIdx != -1 {
			beforeAt := dsn[idx+3 : idx+3+atIdx]
			if colonIdx := strings.Index(beforeAt, ":"); colonIdx != -1 {
				user := beforeAt[:colonIdx]
				return dsn[:idx+3] + user + ":****@" + dsn[idx+3+atIdx+1:]
			}
		}
	}
	return dsn
}
