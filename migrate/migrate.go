package main

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"strconv"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	sqlFilePath := "file://migration_sql"
	postgresUrl := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		"localhost", // hostOSから実行する想定のため明示的に指定
		"5432",      // host portを明示的に指定
		os.Getenv("POSTGRES_DB"))
	m, err := migrate.New(sqlFilePath, postgresUrl)
	if err != nil {
		log.Fatal(err)
	}

	defer m.Close()

	direction := os.Getenv("MIGRATION_DIRECTION")

	switch direction {
	case "up":
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			log.Fatal(err)
		}
		log.Printf("Create table!!!")
	case "down":
		if err := m.Down(); err != nil && err != migrate.ErrNoChange {
			log.Fatal(err)
		}
		log.Printf("Delete table!!!")
	case "force":
		version := os.Getenv("MIGRATION_VERSION")
		v, err := strconv.Atoi(version)
		if err != nil {
			log.Fatal(err)
		}
		if err := m.Force(v); err != nil {
			log.Fatal(err)
		}
		log.Printf("Force set version to %d", v)
	default:
		log.Fatalf("Unknown migration direction: %s", direction)
	}
}

func P(t interface{}) {
	fmt.Println(reflect.TypeOf(t))
}
