package main

import "net/http"

type apiConfig struct{}

type HealthResponse struct {
	Status string `json:"status"`
}

func (apiCfg *apiConfig) healthz(w http.ResponseWriter, r *http.Request) {
	payload := HealthResponse {
		Status: "ok",
	}

	respondWithJSON(w, http.StatusOK, payload)
}
