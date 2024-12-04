package routes

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sort"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/lynxbites/musiclib"
	"github.com/lynxbites/musiclib/internal/db"
)

func NewRouter() *chi.Mux {

	router := chi.NewRouter()

	router.Route("/api/v1/songs", func(r chi.Router) {
		r.Get("/", getSongList)
		r.Put("/", addSong)
		r.Get("/{songId}", getSong)
		r.Patch("/{songId}", patchSong)
	})

	router.Post("/api/v1/song", func(w http.ResponseWriter, r *http.Request) {
		var bytes []byte
		var requestObject musiclib.Song
		bytes, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "No content: Invalid request body", 400)
		}

		if !json.Valid(bytes) {
			http.Error(w, "No content: Invalid request body", 400)
		}
		err = json.Unmarshal(bytes, &requestObject)
		if err != nil {
			http.Error(w, "No content", 400)
		}
		w.Write(bytes)
		fmt.Printf("requestObject: %+v\n", requestObject)

	})
	return router
}

func getSongList(w http.ResponseWriter, r *http.Request) {

	conn := db.New()
	defer conn.Close(context.Background())

	query, err := conn.Query(context.Background(), "select * from songs")
	if err != nil {
		log.Printf("Encountered error when trying to get song list: %v", err)
		http.Error(w, "Encountered Internal Server Error: "+err.Error(), 500)
	}

	var songs []musiclib.Song
	for query.Next() {
		var song musiclib.Song
		err := query.Scan(&song.Id, &song.Group, &song.Name, &song.ReleaseDate, &song.Lyrics, &song.Link)
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
		case "id":
			sort.Slice(songs, func(i, j int) bool {
				return songs[i].Id < songs[j].Id
			})
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
}

func getSong(w http.ResponseWriter, r *http.Request) {

	id := chi.URLParam(r, "songId")
	conn := db.New()
	defer conn.Close(context.Background())
	var song musiclib.Song

	querySelect, err := conn.Query(context.Background(), "select * from songs where songId = "+id)
	if err != nil {
		log.Printf("Encountered error when trying to get song list: %v", err)
		http.Error(w, "Encountered Internal Server Error: "+err.Error(), 500)
	}
	for querySelect.Next() {
		querySelect.Scan(&song.Id, &song.Group, &song.Name, &song.ReleaseDate, &song.Lyrics, &song.Link)
	}
	if song.Id == "" {
		w.WriteHeader(404)
		w.Write([]byte("404 Not found"))
		return
	}
	songJson, err := json.Marshal(song)
	w.Write(songJson)
}

func addSong(w http.ResponseWriter, r *http.Request) {

}

func patchSong(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "songId")
	fmt.Printf("id: %v\n", id)
	w.Write([]byte(id))
}

func paginate(songs []musiclib.Song, pagesStr string, itemsStr string) ([]musiclib.Song, error) {
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
