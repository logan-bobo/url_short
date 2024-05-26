package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
	serverPort := os.Getenv("SERVER_PORT")

	mux := http.NewServeMux()

	server := &http.Server{
		Addr:    ":" + serverPort,
		Handler: mux,
	}

	apiCfg := apiConfig{}

	mux.HandleFunc("GET /api/v1/healthz", apiCfg.healthz)

	log.Printf("Serving port : %v \n", serverPort)
	log.Fatal(server.ListenAndServe())
}
