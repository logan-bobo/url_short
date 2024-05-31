package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"url-short/internal/database"

	_ "github.com/lib/pq"
)

func main() {
	serverPort := os.Getenv("SERVER_PORT")
	dbURL := os.Getenv("PG_CONN")

	db, err := sql.Open("postgres", dbURL)

	if err != nil {
		log.Fatal(err)
	}

	dbQueries := database.New(db)

	mux := http.NewServeMux()

	server := &http.Server{
		Addr:    ":" + serverPort,
		Handler: mux,
	}

	apiCfg := apiConfig{
		DB: dbQueries,
	}

	mux.HandleFunc("GET /api/v1/healthz", apiCfg.healthz)
	mux.HandleFunc("POST /api/v1/data/shorten", apiCfg.postLongURL)

	log.Printf("Serving port : %v \n", serverPort)
	log.Fatal(server.ListenAndServe())
}
