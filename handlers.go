package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"time"

	"url-short/internal/database"
	"url-short/internal/shortener"
)

type apiConfig struct {
	DB *database.Queries
}

type HealthResponse struct {
	Status string `json:"status"`
}

type POSTLongURLRequest struct {
	LongURL string `json:"long_url"`
}

type POSTLongURLResponse struct {
	ShortURL string `json:"short_url"`
}

func (apiCfg *apiConfig) healthz(w http.ResponseWriter, r *http.Request) {
	payload := HealthResponse{
		Status: "ok",
	}

	respondWithJSON(w, http.StatusOK, payload)
}

func (apiCfg *apiConfig) postLongURL(w http.ResponseWriter, r *http.Request) {
	payload := POSTLongURLRequest{}

	err := json.NewDecoder(r.Body).Decode(&payload)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "incorrect request fromat")
		return
	}

	url, err := url.ParseRequestURI(payload.LongURL)

	if err != nil {
		log.Println(err)
		respondWithError(w, http.StatusBadRequest, "could not parse request URL")
		return
	}

	shortURLHash, err := hashCollisionDetection(apiCfg.DB, url.String(), 1, r.Context())

	if err != nil {
		log.Println(err)
		respondWithError(w, http.StatusInternalServerError, "could not resolve hash collision")
	}

	now := time.Now()
	shortenedURL, err := apiCfg.DB.CreateURL(r.Context(), database.CreateURLParams{
		LongUrl:   url.String(),
		ShortUrl:  shortURLHash,
		CreatedAt: now,
		UpdatedAt: now,
	})

	if err != nil {
		log.Println(err)
		respondWithError(w, http.StatusInternalServerError, "could not create short URL in database")
		return
	}

	respondWithJSON(w, http.StatusCreated, POSTLongURLResponse{
		ShortURL: shortenedURL.ShortUrl,
	})
}

func hashCollisionDetection(DB *database.Queries, url string, count int, requestContext context.Context) (string, error) {
	hashURL := shortener.Hash(url, count)
	shortURLHash := shortener.Shorten(hashURL)

	_, err := DB.SelectURL(requestContext, shortURLHash)

	if err == sql.ErrNoRows {
		return shortURLHash, nil
	}

	if err != nil && err != sql.ErrNoRows {
		return "", err
	}

	count++

	return hashCollisionDetection(DB, url, count, requestContext)
}

func (apiCfg *apiConfig) getShortURL(w http.ResponseWriter, r *http.Request) {
	query := r.PathValue("shortUrl")

	row, err := apiCfg.DB.SelectURL(r.Context(), query)

	if err != nil {
		log.Println(err)
		respondWithError(w, http.StatusInternalServerError, "database error")
		return
	}

	http.Redirect(w, r, row.LongUrl, http.StatusMovedPermanently)
}
