package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sort"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/joho/godotenv"
)

type Song struct {
	Group       string      `json:"group"`
	Name        string      `json:"name"`
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

	//Insertions for testing
	_, err = conn.Exec(context.Background(), "INSERT INTO songs VALUES ('C group', 'A name', '2000-01-01', 'text', 'link')")
	if err != nil {
		log.Fatal(err)
	}

	_, err = conn.Exec(context.Background(), "INSERT INTO songs VALUES ('B group', 'C name', '2000-01-03', 'text', 'link')")
	if err != nil {
		log.Fatal(err)
	}

	_, err = conn.Exec(context.Background(), "INSERT INTO songs VALUES ('A group', 'B name', '2000-01-02', 'text', 'link')")
	if err != nil {
		log.Fatal(err)
	}

	_, err = conn.Exec(context.Background(), "INSERT INTO songs VALUES ('D group', 'E name', '2000-01-05', 'text', 'link')")
	if err != nil {
		log.Fatal(err)
	}

	_, err = conn.Exec(context.Background(), "INSERT INTO songs VALUES ('F group', 'D name', '2000-01-06', 'text', 'link')")
	if err != nil {
		log.Fatal(err)
	}

	_, err = conn.Exec(context.Background(), "INSERT INTO songs VALUES ('E group', 'F name', '2000-01-04', 'text', 'link')")
	if err != nil {
		log.Fatal(err)
	}
	//////////////////////////////////////////////////

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
			err := query.Scan(&song.Group, &song.Name, &song.ReleaseDate, &song.Lyrics, &song.Link)
			if err != nil {
				log.Printf("Encountered error when scanning row: %v", err)
				http.Error(w, "Encountered Internal Server Error: "+err.Error(), 500)
				break
			}
			songs = append(songs, song)
		}

		paramSort := r.URL.Query().Get("sort")
		paramPage := r.URL.Query().Get("page")
		paramItems := r.URL.Query().Get("items")

		if paramSort != "" {

			switch paramSort {
			case "group":
				sort.Slice(songs, func(i, j int) bool {
					return songs[i].Group < songs[j].Group
				})
			case "name":
				sort.Slice(songs, func(i, j int) bool {
					return songs[i].Name < songs[j].Name
				})
			case "date":
				sort.Slice(songs, func(i, j int) bool {
					return songs[i].ReleaseDate.Time.Unix() < songs[j].ReleaseDate.Time.Unix()
				})
			case "lyrics":
				sort.Slice(songs, func(i, j int) bool {
					return songs[i].Lyrics < songs[j].Lyrics
				})
			case "link":
				sort.Slice(songs, func(i, j int) bool {
					return songs[i].Link < songs[j].Link
				})
			}
		}

		songs, err = paginate(songs, paramPage, paramItems)

		jsonResp, err := json.MarshalIndent(songs, "", " ")
		if err != nil {
			log.Printf("Encountered error when scanning row: %v", err)
			http.Error(w, "Encountered Internal Server Error: "+err.Error(), 500)
		}
		w.Write(jsonResp)
	})

	http.ListenAndServe(":8000", router)
}

func paginate(songs []Song, pagesStr string, itemsStr string) ([]Song, error) {
	page := 1
	items := 2
	var err error
	if pagesStr != "" {
		page, err = strconv.Atoi(pagesStr)
		if err != nil {
			page = 1
		}
	}

	if itemsStr != "" {
		items, err = strconv.Atoi(itemsStr)
		if err != nil {
			items = 2
		}
	}

	offset := (items * page) - items
	limit := offset + items
	if limit > len(songs) {
		limit = len(songs)
	}
	if offset > len(songs) {
		offset = len(songs)
	}
	songs = songs[offset:limit]

	return songs, nil
}
