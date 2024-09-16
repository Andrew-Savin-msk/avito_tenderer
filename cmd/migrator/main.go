package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	var migrationsPath, direction string

	flag.StringVar(&migrationsPath, "migrations-path", "", "path to migrations")
	flag.StringVar(&direction, "direction", "up", "name of migrations table")
	flag.Parse()

	if migrationsPath == "" {
		log.Fatal("migrations-path is required")
	}

	if migrationsPath == "" {
		log.Fatal("migrations-path is required")
	}

	connPath, exists := os.LookupEnv("POSTGRES_CONN")
	if !exists {
		log.Fatal("no path for migraions")
	}

	m, err := migrate.New("file://"+migrationsPath, fmt.Sprintf("%s", connPath))
	if err != nil {
		log.Fatal(err)
	}

	switch direction {
	case "up":
		err = m.Up()

	case "down":
		err = m.Down()
	}
	if err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			fmt.Println("no changes to apply")

			return
		}
		log.Fatal(err)
	}

	fmt.Println("migrations applied successfuly")
}
