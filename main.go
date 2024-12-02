package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/joho/godotenv"
)

type Song struct {
	group       string
	name        string
	releaseDate pgtype.Date
	lyrics      string
	link        string
}

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

	_, err = conn.Exec(context.Background(), "INSERT INTO songs VALUES ('group', 'name', '2000-01-01', 'text', 'link')")
	if err != nil {
		log.Fatal(err)
	}

	query, err := conn.Query(context.Background(), "SELECT * FROM public.songs")
	if err != nil {
		log.Fatal(err)
	}

	var song Song
	for query.Next() {
		err := query.Scan(&song.group, &song.name, &song.releaseDate, &song.lyrics, &song.link)
		if err != nil {
			log.Printf("Encountered error when scanning row: %v", err)
			break
		}
	}
	fmt.Printf("song: %v\n", song)
	fmt.Printf("time: %v\n", song.releaseDate.Time.UTC())
}
