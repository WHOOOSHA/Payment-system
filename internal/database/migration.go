package database

import (
	"fmt"
	"log"
	"payment_system/configs"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func RunMigrations(config *configs.DatabaseConfig) {
	connStr := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		config.User,
		config.Password,
		config.Host,
		config.Port,
		config.DBname,
		config.SSLmode,
	)

	m, err := migrate.New(
		"file://migrations",
		connStr,
	)
	if err != nil {
		log.Fatal("Migration failed: ", err)
	}

	if err := m.Up(); err != nil && err.Error() != "no change" {
		log.Fatal("Migration failed: ", err)
	}
}
