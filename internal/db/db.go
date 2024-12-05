package db

import (
	"context"
	"fmt"
	"os"

	"github.com/charmbracelet/log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

type DB struct {
	*pgx.Conn
}

var connstr string
var connstrm string

func init() {
	log.SetReportCaller(true)

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}
	var connstrfound bool
	connstr, connstrfound = os.LookupEnv("CONNSTR")
	if connstr == "" || connstrfound == false {
		log.Fatal("lookup connection string: connection string not specified in .env file")
	}
	connstrm, connstrfound = os.LookupEnv("CONNSTRMIGRATION")
	if connstr == "" || connstrfound == false {
		log.Fatal("lookup connection string: connection string not specified in .env file")
	}
}

func New() *DB {

	conn, err := pgx.Connect(context.Background(), connstr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	return &DB{
		conn,
	}
}

func Migration() (*migrate.Migrate, error) {

	m, err := migrate.New("file://internal/db/migrations/", connstrm)
	if err != nil {
		return nil, err
	}
	return m, nil
}
