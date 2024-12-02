package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	connstr, connstrfound := os.LookupEnv("CONNSTR")
	if connstr == "" || connstrfound == false {
		log.Fatalf("lookup connection string: connection string not specified in .env file")
	}

	conn, err := pgx.Connect(context.Background(), connstr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	conn.Exec(context.Background(), "INSERT INTO songs (Group, Name, Lyrics, Link, Date ) VALUES ('group', 'name', 'text', 'https', '2000-01-01')")

}
