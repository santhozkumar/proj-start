package database

import (
	"database/sql"
	"log"
	db "osvauld/db/sqlc"

	_ "github.com/lib/pq"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

var Store db.Store

func DBConnection(masterDSN string) error {
	conn, err := sql.Open("postgres", masterDSN)
	// db, err := sql.Open("postgres", "postgres://localhost:5432/database?sslmode=enable")
	driver, err := postgres.WithInstance(conn, &postgres.Config{})
	m, err := migrate.NewWithDatabaseInstance(
		"file://db/migration",
		"postgres", driver)
	m.Up()

	if err != nil {
		log.Fatalln("Connection faild")
		return err
	}
	if err := conn.Ping(); err != nil {
		log.Fatalln("Ping failed")
		return err
	}

	Store = db.NewStore(conn)
	return nil
}
