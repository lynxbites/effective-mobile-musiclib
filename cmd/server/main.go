package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/lynxbites/musiclib/internal/db"
	"github.com/lynxbites/musiclib/internal/routes"
	_ "github.com/swaggo/http-swagger/example/go-chi/docs"
	_ "github.com/swaggo/http-swagger/v2"
)

// @title MusicLib
// @version 0.1
// @description This is a simple swagger for musiclib.

// @host localhost:8000
// @BasePath /api/
func init() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}

	mydir, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(mydir)
}

func main() {

	router := routes.NewRouter()

	m, err := db.Migration()
	if err != nil {
		log.Fatal(err)
	}
	m.Up()

	http.ListenAndServe(":8000", router)

}