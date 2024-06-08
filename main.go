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
	jwtSecret := os.Getenv("JWT_SECRET")

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
		DB:        dbQueries,
		JWTSecret: jwtSecret,
	}

	// utility endpoints
	mux.HandleFunc("GET /api/v1/healthz", apiCfg.healthz)

	// url management endpoints
	mux.HandleFunc("POST /api/v1/data/shorten", apiCfg.postLongURL)
	mux.HandleFunc("GET /api/v1/{shortUrl}", apiCfg.getShortURL)

	// user management endpoints
	mux.HandleFunc("POST /api/v1/users", apiCfg.postAPIUsers)
	mux.HandleFunc("PUT /api/v1/users", apiCfg.authenticationMiddlewear(apiCfg.putAPIUsers))
	mux.HandleFunc("POST /api/v1/login", apiCfg.postAPILogin)

	log.Printf("Serving port : %v \n", serverPort)
	log.Fatal(server.ListenAndServe())
}
