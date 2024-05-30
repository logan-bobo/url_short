package main

import (
	"encoding/json"
	"net/http"
	"net/url"
	"url-short/internal/database"
)

type apiConfig struct{
	DB *database.Queries
}

type HealthResponse struct {
	Status string `json:"status"`
}

type POSTLongURLRequest struct {
	LongURL string `json:"status"`
}

func (apiCfg *apiConfig) healthz(w http.ResponseWriter, r *http.Request) {
	payload := HealthResponse {
		Status: "ok",
	}

	respondWithJSON(w, http.StatusOK, payload)
}

func (apiCfg *apiConfig) POSTLongURL(w http.ResponseWriter, r *http.Request) {
	payload := POSTLongURLRequest{}

	err := json.NewDecoder(r.Body).Decode(&payload)

	if err != nil {
		respondWithError(w, 400, "incorrect request fromat")
		return
	}

	url, err := url.ParseRequestURI(payload.LongURL)

	if err != nil {
		respondWithError(w, 400, "could not parse request URL")
		return
	}
	
	shortURL, err :=apiCfg.DB.CreateURL()	
}
