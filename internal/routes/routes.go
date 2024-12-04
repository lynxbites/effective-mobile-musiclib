package routes

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"sort"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/lynxbites/musiclib"
	"github.com/lynxbites/musiclib/internal/db"
	_ "github.com/swaggo/http-swagger/example/go-chi/docs"
)

func NewRouter() *chi.Mux {

	router := chi.NewRouter()

	router.Route("/api/v1/songs", func(r chi.Router) {
		r.Get("/", getSongList)
		r.Post("/", addSong)
		r.Get("/{songId}", getSongs)
		r.Patch("/{songId}", patchSong)
		r.Delete("/{songId}", deleteSong)
	})

	return router
}

// ListAccounts godoc
// @Summary      Get song
// @Description  Gets song from db
// @Tags         Songs
// @Param   sort      query     string     false  "string valid" minValue(0)
// @Param   page      query     int     false  "int > 0"
// @Param   items      query     int     false  "int > 0"
// @Produce      json
// @Success      200 {object} []musiclib.Song "OK"
// @Failure      400  "Bad Request"
// @Failure      500  "Internal error"
// @Router       /v1/songs [get]
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
		err := query.Scan(&song.Id, &song.Group, &song.Name, &song.ReleaseDate, &song.Text, &song.Link)
		if err != nil {
			log.Printf("Encountered error when scanning row: %v", err)
			http.Error(w, "Encountered Internal Server Error: "+err.Error(), 500)
			return
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
		case "text":
			sort.Slice(songs, func(i, j int) bool {
				return songs[i].Text < songs[j].Text
			})
		case "link":
			sort.Slice(songs, func(i, j int) bool {
				return songs[i].Link < songs[j].Link
			})
		}
	}

	page := 1
	items := 2

	if paramPage != "" {
		page, err = strconv.Atoi(paramPage)
		if err != nil {
			page = 1
		}
		if page <= 0 {
			http.Error(w, "Bad Request", 400)
			return
		}

	}

	if paramItems != "" {
		items, err = strconv.Atoi(paramItems)
		if err != nil {
			items = 2
		}
		if items <= 0 {
			http.Error(w, "Bad Request", 400)
			return
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

	jsonResp, err := json.MarshalIndent(songs, "", " ")
	if err != nil {
		log.Printf("Encountered error when scanning row: %v", err)
		http.Error(w, "Encountered Internal Server Error: "+err.Error(), 500)
		return
	}
	w.Write(jsonResp)
}

// ListAccounts godoc
// @Summary      Get song
// @Description  Gets song from db
// @Tags         Songs
// @Produce      json
// @Success      200 {array} musiclib.Song "OK"
// @Failure      400  "Bad Request"
// @Failure      404  "Not Found"
// @Failure      500  "Internal error"
// @Router       /v1/songs/{songId} [get]
func getSongs(w http.ResponseWriter, r *http.Request) {

	id := chi.URLParam(r, "songId")
	conn := db.New()
	defer conn.Close(context.Background())
	var song musiclib.Song

	querySelect, err := conn.Query(context.Background(), "select * from songs where songId = "+id)
	if err != nil {
		log.Printf("Encountered error when trying to get song list: %v", err)
		http.Error(w, "Encountered Internal Server Error: "+err.Error(), 500)
		return
	}
	for querySelect.Next() {
		err := querySelect.Scan(&song.Id, &song.Group, &song.Name, &song.ReleaseDate, &song.Text, &song.Link)
		if err != nil {
			log.Printf("Encountered error when trying to scan query rows: %v", err)
			http.Error(w, "Encountered Internal Server Error: "+err.Error(), 500)
			return
		}
	}

	if song.Id == "" {
		http.Error(w, "404 Not found", 404)
		return
	}
	songJson, err := json.Marshal(song)
	w.Write(songJson)
}

// ListAccounts godoc
// @Summary      Post song
// @Description  Posts song to db
// @Tags         Songs
// @Accept       json
// @Param 		 json body string true "Song JSON Object" SchemaExample({\n\t"group":"Johnny Mercer",\n\t"name":"Personality",\n\t"releaseDate":"1946-12-14",\n\t"text":"When Madam Pompadour was on a ballroom floor\n\tSaid all the gentlemen, "Obviously"\n\t"The madam has the cutest personality"\n\tAnd think of all the books about do Barry's looks\n\tWhat was it made her the toast of Paree?\n\tShe had a well-developed personality",\n\t"https://www.youtube.com/watch?v=c1L6ZdbL5a0":"rer"\n})
// @Produce      json
// @Success      200  "OK"
// @Failure      400  "Bad Request"
// @Failure      409  "Conflict"
// @Failure      500  "Internal error"
// @Router       /v1/songs [post]
func addSong(w http.ResponseWriter, r *http.Request) {

	conn := db.New()
	defer conn.Close(context.Background())

	var song musiclib.SongPost
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&song)
	if err != nil {
		log.Printf("POST /songs: Error while decoding JSON: %+v\n", err)
		http.Error(w, "Invalid JSON data", 400)
		return
	}
	fmt.Printf("song: %+v\n", song)
	if !isPostValid(song) {
		log.Printf("POST /songs: Invalid JSON data: %+v\n", err)
		http.Error(w, "Invalid JSON data", 400)
		return
	}
	if decoder.More() {
		log.Printf("POST /songs: Additional data")
		http.Error(w, "Additional data", 400)
		return
	}

	var exists bool
	queryRows := conn.QueryRow(context.Background(), "select exists(select 1 from songs where groupName=$1 and songName=$2)", song.Group, song.Name)
	queryRows.Scan(&exists)
	if exists {
		log.Printf("POST /songs: Song already exists")
		http.Error(w, "Song already exists", 409)
		return
	}
	_, err = conn.Exec(context.Background(), "insert into songs (groupName, songName, releaseDate, songText, songLink) values ('$1','$2','$3','$4','$5')", song.Group, song.Name, song.ReleaseDate, song.Text, song.Link)
	if err != nil {
		log.Printf("Encountered error when trying to insert song data: %v", err)
		http.Error(w, "Encountered Internal Server Error: "+err.Error(), 500)
		return
	}
	w.WriteHeader(200)

}

func postSong(w http.ResponseWriter, r *http.Request) {

}

func patchSong(w http.ResponseWriter, r *http.Request) {

}

// ListAccounts godoc
// @Summary      Delete song
// @Description  deletes song from db
// @Tags         Songs
// @Produce      json
// @Success      200,204 "OK"
// @Failure      400  "Bad Request"
// @Failure      404  "Not Found"
// @Failure      500  "Internal error"
// @Router       /v1/songs/{songId} [delete]
func deleteSong(w http.ResponseWriter, r *http.Request) {

	id := chi.URLParam(r, "songId")
	conn := db.New()
	defer conn.Close(context.Background())

	exists, err := isRowExists(conn, "songId", id)
	if err != nil {
		log.Print(err)
		http.Error(w, "Encountered Internal Server Error: "+err.Error(), 500)
		return
	}
	_, err = conn.Exec(context.Background(), "delete from songs where songId = $1", id)
	if err != nil {
		log.Printf("Encountered error when trying to delete data: %v", err)

		return
	}

	if exists {
		w.WriteHeader(200)
	} else {
		w.WriteHeader(204)
	}

}

func isRowExists(conn *db.DB, key string, value string) (bool, error) {
	var exists bool

	queryRow := conn.QueryRow(context.Background(), "select exists(select 1 from songs where $1=$2)", key, value)
	err := queryRow.Scan(&exists)
	if err != nil {
		return false, errors.New("Encountered error when trying to check if row exists: " + err.Error())
	}
	return exists, nil
}

func isPostValid(patch musiclib.SongPost) bool {
	if patch.Group == nil {
		log.Print("Invalid Group")
		return false
	}
	if patch.Name == nil {
		log.Print("Invalid Name")
		return false
	}
	if patch.ReleaseDate == nil {
		log.Print("Invalid releaseDate")
		return false
	}
	if patch.Text == nil {
		log.Print("Invalid Text")
		return false
	}
	if patch.Link == nil {
		log.Print("Invalid Link")
		return false
	}

	return true
}
