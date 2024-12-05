package main

import (
	"net/http"

	"github.com/charmbracelet/log"
	"github.com/joho/godotenv"
	"github.com/lynxbites/musiclib/internal/db"
	"github.com/lynxbites/musiclib/internal/routes"
	_ "github.com/swaggo/http-swagger/example/go-chi/docs"
	_ "github.com/swaggo/http-swagger/v2"
)

// @title MusicLib
// @version 0.3
// @description This is a simple swagger for musiclib.
// @schemes http
// @host localhost:8000
// @BasePath /api/
var runSwagger bool

func init() {
	log.SetReportCaller(true)

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	routerSwagger := routes.NewSwaggerRouter()

	router := routes.NewRouter()

	m, err := db.Migration()
	if err != nil {
		log.Fatal(err)
	}
	m.Up()

	log.Info("Starting Swagger server...")
	go http.ListenAndServe(":8001", routerSwagger)
	log.Info("Starting API...")
	http.ListenAndServe(":8000", router)

}
