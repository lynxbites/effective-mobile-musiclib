package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/lynxbites/musiclib/internal/db"
	"github.com/lynxbites/musiclib/internal/routes"
	flag "github.com/spf13/pflag"
	_ "github.com/swaggo/http-swagger/example/go-chi/docs"
	_ "github.com/swaggo/http-swagger/v2"
)

// @title MusicLib
// @version 0.1
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

	flag.BoolVarP(&runSwagger, "swagger", "s", false, "Set to run swagger server.")
	flag.Lookup("swagger").NoOptDefVal = "true"

	mydir, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(mydir)
}

func main() {
	routerSwagger := routes.NewSwaggerRouter()

	router := routes.NewRouter()

	m, err := db.Migration()
	if err != nil {
		log.Fatal(err)
	}
	m.Up()

	log.Printf("Starting Swagger server...")
	go http.ListenAndServe(":8001", routerSwagger)
	log.Printf("Starting API...")
	http.ListenAndServe(":8000", router)

}
