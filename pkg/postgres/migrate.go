package postgres

import (
	"dynamic-user-segmentation/pkg/logging"
	"errors"
	"github.com/golang-migrate/migrate/v4"
	"time"
)

const (
	maxDefaultMigrationAttempts = 5
	defaultTimeout              = time.Second
)

func runMigrations(dsn string, log logging.Logger) error {
	log.Info("Migrate: starting...")
	dsn += "?sslmode=disable"
	attemptsCount := maxDefaultMigrationAttempts
	var (
		m   *migrate.Migrate
		err error
	)

	for attemptsCount > 0 {
		m, err = migrate.New("file://migrations", dsn)
		if err == nil {
			break
		}
		log.Infof("try to migrate... attempts left %d\n", attemptsCount)
		time.Sleep(defaultTimeout)
		attemptsCount--
	}
	if err != nil {
		return err
	}

	defer func() {
		_, _ = m.Close()
	}()
	if err = m.Up(); !errors.Is(err, migrate.ErrNoChange) {
		return err
	}
	if errors.Is(err, migrate.ErrNoChange) {
		log.Info("Migrate: no change")
		return nil
	}
	log.Info("Migrate: success")
	return nil
}
