package db

import (
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type Migration struct {
	Url           string
	MigrationsDir string
	SourceDriver  source.Driver
}

func NewMigration(url string, migrationsDir string, sourceDriver source.Driver) *Migration {
	return &Migration{
		Url:           url,
		MigrationsDir: migrationsDir,
		SourceDriver:  sourceDriver,
	}
}

func (m *Migration) Migrate() error {
	var migr *migrate.Migrate
	var err error

	if m.SourceDriver != nil {
		migr, err = migrate.NewWithSourceInstance("httpfs", m.SourceDriver, m.Url)
	} else {
		migr, err = migrate.New(fmt.Sprintf("file://%s", m.MigrationsDir), m.Url)
	}

	if err != nil {
		return err
	}
	defer migr.Close()

	err = migr.Up()
	if !errors.Is(err, migrate.ErrNoChange) && err != nil {
		return err
	}
	return nil
}

func (m *Migration) Rollback() error {
	var migr *migrate.Migrate
	var err error

	if m.SourceDriver != nil {
		migr, err = migrate.NewWithSourceInstance("httpfs", m.SourceDriver, m.Url)
	} else {
		migr, err = migrate.New(fmt.Sprintf("file://%s", m.MigrationsDir), m.Url)
	}

	if err != nil {
		return err
	}
	defer migr.Close()

	err = migr.Steps(-1)
	if !errors.Is(err, migrate.ErrNoChange) && err != nil {
		return err
	}
	return nil
}
