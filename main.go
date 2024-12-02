package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	connstr, connstrfound := os.LookupEnv("CONNSTR")
	if connstr == "" || connstrfound == false {
		log.Fatalf("lookup connection string: connection string not specified in .env file")
	}
}
