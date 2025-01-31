package main

import (
	"errors"
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/juliendoutre/recorder/internal/config"
	"go.uber.org/zap"
)

//nolint:gochecknoglobals
var (
	GoVersion string
	Os        string //nolint:varnamelen
	Arch      string
)

func main() {
	logger, err := zap.NewProductionConfig().Build()
	if err != nil {
		log.Panic(err)
	}

	defer func() { _ = logger.Sync() }()

	pgURL, err := config.PostgresURL()
	if err != nil {
		logger.Panic("Reading PostgresQL config", zap.Error(err))
	}

	migrationURL := config.MigrationsURL()

	migrator, err := migrate.New(migrationURL.String(), pgURL.String())
	if err != nil {
		logger.Panic("Creating migrator", zap.Error(err))
	}

	if err := migrator.Up(); err != nil {
		if !errors.Is(err, migrate.ErrNoChange) {
			logger.Panic("Running migrations", zap.Error(err))
		}
	}
}
