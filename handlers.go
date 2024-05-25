package main

import "net/http"

type apiConfig struct{}

func (apiCfg *apiConfig) healthz(w http.ResponseWriter, r *http.Request) {
	payload := struct {
		Status string `json:"status"`
	}{
		Status: "ok",
	}

	respondWithJSON(w, http.StatusOK, payload)
}
