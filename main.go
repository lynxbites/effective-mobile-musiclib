package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/joho/godotenv"
)

type Song struct {
	Group       string      `json:"group"`
	Game        string      `json:"name"`
	ReleaseDate pgtype.Date `json:"releaseDate"`
	Lyrics      string      `json:"lyrics"`
	Link        string      `json:"link"`
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

	_, err = conn.Exec(context.Background(), "INSERT INTO songs VALUES ('group', 'name', '2000-01-01', 'text', 'link')")
	if err != nil {
		log.Fatal(err)
	}

	_, err = conn.Exec(context.Background(), "INSERT INTO songs VALUES ('group', 'name', '2000-01-01', 'text', 'link')")
	if err != nil {
		log.Fatal(err)
	}

	// fmt.Printf("song: %v\n", song)
	// fmt.Printf("time: %v\n", song.releaseDate.Time.UTC())

	router := chi.NewRouter()

	router.Get("/api/v1/list", func(w http.ResponseWriter, r *http.Request) {

		query, err := conn.Query(context.Background(), "SELECT * FROM public.songs")
		if err != nil {
			log.Printf("Encountered error when trying to get song list: %v", err)
			http.Error(w, "Encountered Internal Server Error: "+err.Error(), 500)
		}

		var songs []Song
		for query.Next() {
			var song Song
			err := query.Scan(&song.Group, &song.Game, &song.ReleaseDate, &song.Lyrics, &song.Link)
			if err != nil {
				log.Printf("Encountered error when scanning row: %v", err)
				http.Error(w, "Encountered Internal Server Error: "+err.Error(), 500)
				break
			}
			songs = append(songs, song)
		}
		jsonResp, err := json.MarshalIndent(songs, "", " ")
		if err != nil {
			log.Printf("Encountered error when scanning row: %v", err)
			http.Error(w, "Encountered Internal Server Error: "+err.Error(), 500)
		}
		w.Write(jsonResp)
	})

	http.ListenAndServe(":8000", router)
}
