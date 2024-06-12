package main

import (
	"database/sql"
	"log"
	"os"
	"picturethisai/db"

	"github.com/joho/godotenv"
	"github.com/mattes/migrate"
	"github.com/mattes/migrate/database/postgres"
)

func createDB() (*sql.DB, error) {
	if err := godotenv.Load(); err != nil {
		return nil, err
	}
	var (
		host   = os.Getenv("DB_HOST")
		user   = os.Getenv("DB_USER")
		pass   = os.Getenv("DB_PASSWORD")
		dbname = os.Getenv("DB_NAME")
	)
	return db.CreateDatabase(dbname, user, pass, host)
}

func main() {
	db, err := createDB()
	if err != nil {
		log.Fatal(err)
	}

	// Create migration instance
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatal(err)
	}

	// Point to your migration files. Here we're using local files, but it could be other sources.
	m, err := migrate.NewWithDatabaseInstance(
		"file://cmd/migrate/migrations", // source URL
		"postgres",                      // database name
		driver,                          // database instance
	)
	if err != nil {
		log.Fatal(err)
	}

	cmd := os.Args[len(os.Args)-1]
	if cmd == "up" {
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			log.Fatal(err)
		}
	}
	if cmd == "down" {
		if err := m.Down(); err != nil && err != migrate.ErrNoChange {
			log.Fatal(err)
		}
	}
}
