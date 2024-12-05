package routes

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/lynxbites/musiclib"
	"github.com/lynxbites/musiclib/internal/db"
	_ "github.com/swaggo/http-swagger/example/go-chi/docs"
)

func NewRouter() *chi.Mux {

	router := chi.NewRouter()

	router.Route("/api/v1/songs", func(r chi.Router) {
		r.Use(middleware.Logger)
		r.Use(cors.Handler(cors.Options{
			AllowedOrigins:   []string{"http://*"},
			AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE"},
			AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
			ExposedHeaders:   []string{"Link"},
			AllowCredentials: false,
			MaxAge:           360,
		}))
		r.Get("/", getSongList)
		r.Post("/", addSong)
		r.Get("/{songId}", getSong)
		r.Patch("/{songId}", patchSong)
		r.Delete("/{songId}", deleteSong)
	})

	return router
}

// GetSongList godoc
// @Summary      Get songs
// @Description  Gets list of songs from DB, with filters and pagination.
// @Tags         Songs
// @Param   filter      query     string     false  "Filter by id, group, name, date, text or link."
// @Param   page      query     int     false 	"Number of the page."
// @Param   items      query     int     false 	"How many items to display per page."
// @Produce      json
// @Success      200 {array} musiclib.Song "OK"
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

	paramFilter := r.URL.Query().Get("filter")
	paramPage := r.URL.Query().Get("page")
	paramItems := r.URL.Query().Get("items")

	if paramFilter != "" {

		switch paramFilter {
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
				return songs[i].ReleaseDate < songs[j].ReleaseDate
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

	encoder := json.NewEncoder(w)
	encoder.Encode(songs)
}

// GetSong godoc
// @Summary      Get song
// @Description  Get a song from DB, with pagination for verses.
// @Tags         Songs
// @Produce      json
// @Param   offset      query     int     false 	"Verse offset."
// @Param   limit      query     int     false		"How many verses to display."
// @Param   	 songId      path     int     true  "Id of the song."
// @Success      200 {object} musiclib.SongPaginated "OK"
// @Failure      400  "Bad Request"
// @Failure      404  "Not Found"
// @Failure      500  "Internal error"
// @Router       /v1/songs/{songId} [get]
func getSong(w http.ResponseWriter, r *http.Request) {
	paramId := chi.URLParam(r, "songId")
	paramOffset := r.URL.Query().Get("offset")
	paramLimit := r.URL.Query().Get("limit")
	conn := db.New()
	defer conn.Close(context.Background())
	var song musiclib.Song

	querySelect, err := conn.Query(context.Background(), "select * from songs where songId = "+paramId)
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

	presplit := song.Text
	splitter := `\n`
	textParsed := strings.Split(presplit, splitter)
	textParsed = removeEmptyStrings(textParsed)
	fmt.Printf("len(splitted): %v\n", len(textParsed))

	songPaginated := musiclib.SongPaginated{
		Id:          song.Id,
		Group:       song.Group,
		Name:        song.Name,
		ReleaseDate: song.ReleaseDate,
		Text:        textParsed,
		Link:        song.Link,
	}

	offset := 0
	limit := 3

	if paramOffset != "" {
		offset, err = strconv.Atoi(paramOffset)
		if err != nil {
			offset = 0
		}
		if offset < 0 {
			http.Error(w, "Bad Request", 400)
			return
		}

	}

	if paramLimit != "" {
		limit, err = strconv.Atoi(paramLimit)
		if err != nil {
			limit = 3
		}
		if limit <= 0 {
			http.Error(w, "Bad Request", 400)
			return
		}
	}

	limit = offset + limit

	if limit > len(songPaginated.Text) {
		limit = len(songPaginated.Text)
	}
	if offset > len(songPaginated.Text) {
		offset = len(songPaginated.Text)
	}
	songPaginated.Text = songPaginated.Text[offset:limit]

	encoder := json.NewEncoder(w)
	encoder.Encode(songPaginated)

}

// AddSong godoc
// @Summary      Post song
// @Description  Post song to DB.
// @Tags         Songs
// @Accept       json
// @Param 		 json body string true "Song JSON Object" SchemaExample({"group":"Author name", "name":"Song name", "releaseDate":"2024-12-12", "text":"Lyrics", "link":"Link"})
// @Produce      json
// @Success      200  "OK"
// @Failure      400  "Bad Request"
// @Failure      409  "Conflict"
// @Failure      500  "Internal error"
// @Router       /v1/songs [post]
func addSong(w http.ResponseWriter, r *http.Request) {

	log.Printf("POST /v1/songs request from %v", r.Host)
	conn := db.New()
	defer conn.Close(context.Background())

	var songPost musiclib.SongPost
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&songPost)
	if err != nil {
		log.Printf("Error while decoding JSON: %+v\n", err)
		http.Error(w, "Invalid JSON data", 400)
		return
	}
	fmt.Printf("song: %+v\n", songPost)
	if !isPostValid(songPost) {
		log.Printf("Invalid JSON data: %+v\n", err)
		http.Error(w, "Invalid JSON data", 400)
		return
	}
	if decoder.More() {
		log.Printf("Additional data")
		http.Error(w, "Additional data", 400)
		return
	}

	var exists bool
	queryRows := conn.QueryRow(context.Background(), "select exists(select 1 from songs where groupName=$1 and songName=$2)", *songPost.Group, *songPost.Name)
	queryRows.Scan(&exists)
	if exists {
		log.Printf("Song already exists")
		http.Error(w, "Song already exists", 409)
		return
	}

	_, err = conn.Exec(context.Background(), `insert into songs (groupName, songName, releaseDate, songText, songLink) values ($1,$2,$3,$4,$5)`, *songPost.Group, *songPost.Name, *songPost.ReleaseDate, *songPost.Text, *songPost.Link)
	if err != nil {
		log.Printf("Encountered error when trying to insert song data: %v", err)
		http.Error(w, "Encountered Internal Server Error: "+err.Error(), 500)
		return
	}

	w.WriteHeader(200)

	//Get /info request // Я так и не понял что от меня требуется во втором задании, извиняюсь за недопонимание :C
	requestString := *songPost.Group + "&" + "name=" + *songPost.Name
	request, err := http.Get("https://example.com/info?group=" + url.QueryEscape(requestString))
	if err != nil {
		log.Printf("Encountered error when requesting /info: %v", err)
	}
	fmt.Printf("request status: %v\n", request.Status)
	body, err := io.ReadAll(request.Body)
	if err != nil {
		log.Printf("/info body read error: %v", err)
	}
	fmt.Printf("request body: %v\n", string(body))
}

// PatchSong godoc
// @Summary      Patch song
// @Description  Update song specified by id.
// @Tags         Songs
// @Produce      json
// @Param 		 json body string true "Song JSON Object" SchemaExample({"group":"Patched", "name":"PatchedName", "releaseDate":"2023-12-12", "text":"PatchedText", "link":"PatchedLink"})
// @Param   	 songId      path     int     true  "Id of a song to patch."
// @Success      200  "OK"
// @Failure      400  "Bad Request"
// @Failure      404  "Not Found"
// @Failure      500  "Internal error"
// @Router       /v1/songs/{songId} [patch]
func patchSong(w http.ResponseWriter, r *http.Request) {
	conn := db.New()
	defer conn.Close(context.Background())

	paramId := chi.URLParam(r, "songId")
	fmt.Printf("paramId: %v\n", paramId)

	var exists bool
	queryRows := conn.QueryRow(context.Background(), "select exists(select 1 from songs where songId = $1)", paramId)
	queryRows.Scan(&exists)
	if !exists {
		log.Printf("Song does not exist")
		http.Error(w, "Song does not exist", 400)
		return
	}

	var patchRequest musiclib.SongPatch
	var patchedObject musiclib.SongPatch
	queryObj := conn.QueryRow(context.Background(), "select groupName, songName, releaseDate, songText, songLink from songs where songId = "+paramId)
	err := queryObj.Scan(&patchedObject.Group, &patchedObject.Name, &patchedObject.ReleaseDate, &patchedObject.Text, &patchedObject.Link)
	if err != nil {
		log.Println("Error while scanning object: ", err)
		http.Error(w, "Error while patching", 500)
		return
	}
	fmt.Printf("patchedObject: %+v\n", patchedObject)
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&patchRequest)
	if err != nil {
		log.Println("Error while decoding body: ", err)
		http.Error(w, "Error while patching", 500)
		return
	}
	fmt.Printf("patchRequest: %+v\n", patchRequest)

	if patchRequest.Group != "" {
		patchedObject.Group = patchRequest.Group
	}
	if patchRequest.Name != "" {
		patchedObject.Name = patchRequest.Name
	}
	if patchRequest.ReleaseDate != "" {
		patchedObject.ReleaseDate = patchRequest.ReleaseDate
	}
	if patchRequest.Text != "" {
		patchedObject.Text = patchRequest.Text
	}
	if patchRequest.Link != "" {
		patchedObject.Link = patchRequest.Link
	}

	fmt.Printf("after patch: %+v\n", patchedObject)

	_, err = conn.Exec(context.Background(), `update songs set groupName = $1, songName = $2, releaseDate = $3, songText = $4, songLink = $5 where songId = $6`, patchedObject.Group, patchedObject.Name, patchedObject.ReleaseDate, patchedObject.Text, patchedObject.Link, paramId)
	if err != nil {
		log.Println("Error while updating patch object: ", err)
		http.Error(w, "Error while patching", 500)
		return
	}
	w.WriteHeader(200)
}

// DeleteSong godoc
// @Summary      Delete song
// @Description  Delete song from DB.
// @Tags         Songs
// @Produce      json
// @Param   	 songId      path     int     true  "Id of a song to delete"
// @Success      200,204 "OK"
// @Failure      400  "Bad Request"
// @Failure      500  "Internal error"
// @Router       /v1/songs/{songId} [delete]
func deleteSong(w http.ResponseWriter, r *http.Request) {

	paramId := chi.URLParam(r, "songId")
	idInt, err := strconv.Atoi(paramId)
	if err != nil {
		log.Printf("Bad request")
		http.Error(w, "Bad request", 400)
	}
	if idInt <= 0 {
		log.Printf("Bad request")
		http.Error(w, "Bad request", 400)
	}

	conn := db.New()
	defer conn.Close(context.Background())

	_, err = conn.Exec(context.Background(), "delete from songs where songId = $1", paramId)
	if err != nil {
		log.Printf("Encountered error when trying to delete data: %v", err)
		http.Error(w, "Error while deleting", 500)
		return
	}

	w.WriteHeader(204)

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

func removeEmptyStrings(arr []string) []string {
	var newArr []string
	for i := range arr {
		if arr[i] != "" {
			strings.TrimSpace(arr[i])
			newArr = append(newArr, arr[i])
		}
	}
	return newArr
}
