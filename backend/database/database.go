package database

import (
	"database/sql"
	_ "embed"
	"events/backend/database/gen"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

//go:generate go run github.com/sqlc-dev/sqlc/cmd/sqlc@latest generate

//go:embed schema.sql
var databaseSchema string

type DB struct {
	*gen.Queries
}

func OpenMemory() *DB {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		log.Fatal(err)
	}

	// initialize the database
	_, err = db.Exec(databaseSchema)
	if err != nil {
		log.Fatal(err)
	}

	return &DB{gen.New(db)}
}

func Open(filename string) *DB {
	db, err := sql.Open("sqlite3", filename)
	if err != nil {
		log.Fatal(err)
	}

	// run migrations
	err = initMigrator().Migrate(db)
	if err != nil {
		log.Fatal("migrations failed:", err)
	}

	return &DB{gen.New(db)}
}

var database *DB = Open("data.db")

func Default() *DB {
	return database
}

func DangerousSwitchDefaultToMemory() {
	database = OpenMemory()
}
