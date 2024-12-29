package database

import (
	"database/sql"
	_ "embed"
	"log"

	"github.com/lopezator/migrator"
)

//go:embed migrations/001_init.sql
var migration001 string

func initMigrator() *migrator.Migrator {
	m, err := migrator.New(
		migrator.Migrations(
			&migrator.Migration{
				Name: "Create initial tables",
				Func: func(tx *sql.Tx) error {
					if _, err := tx.Exec(migration001); err != nil {
						return err
					}
					return nil
				},
			},
		),
	)
	if err != nil {
		log.Fatal(err)
	}

	return m
}
