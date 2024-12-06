package main

import (
	"net/http"
	"os"

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

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}
	log.SetReportCaller(true)
	levelStr, set := os.LookupEnv("LOGLEVEL")
	if !set {
		log.Warn("Log level is not set, defaulting to INFO level.")
	} else {
		level, err := log.ParseLevel(levelStr)
		if err != nil {
			log.Warn("Error setting log level, defaulting to INFO level: " + err.Error())
		} else {
			log.SetLevel(level)
		}

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
