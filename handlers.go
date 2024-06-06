package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/mail"
	"net/url"
	"time"

	"url-short/internal/database"
	"url-short/internal/shortener"

	"golang.org/x/crypto/bcrypt"
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

type POSTAPIUser struct {
	Email string `json:"email"`
	Password string `json:"Password"`
}

type POSTAPIUsersResponse struct {
	ID int32 `json:"id"`
	Email string `json:"email"`
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

func (apiCfg *apiConfig) postAPIUsers(w http.ResponseWriter, r *http.Request) {
	payload := POSTAPIUser{}

	err := json.NewDecoder(r.Body).Decode(&payload)

	if err != nil {
		log.Println(err)
		respondWithError(w, http.StatusBadRequest, "incorrect parameters for user creation")
		return
	}

	_, err = mail.ParseAddress(payload.Email)

	if err != nil {
		fmt.Println(err)
		respondWithError(w, http.StatusBadRequest, "invalid email address")
		return	
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)

	if err != nil {
		log.Println(err)
		respondWithError(w, http.StatusBadRequest, "bad password supplied from client")
		return
	}

	now := time.Now()

	user, err := apiCfg.DB.CreateUser(r.Context(), database.CreateUserParams{
		Email: payload.Email,
		Password: string(passwordHash),
		CreatedAt: now,
		UpdatedAt: now,
	})
	
	respondWithJSON(w, http.StatusCreated, POSTAPIUsersResponse{
		ID: user.ID, 
		Email: user.Email,
	})	
}

func (apiCfg *apiConfig) postAPILogin(w http.ResponseWriter, r *http.Request) {
	payload := POSTAPIUser{}

	err := json.NewDecoder(r.Body).Decode(&payload)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid parameters for login")
		return
	}

	user, err := apiCfg.DB.SelectUser(r.Context(), payload.Email)

	if err == sql.ErrNoRows {
		respondWithError(w, http.StatusNotFound, "could not find user")
		return
	}

	if err != nil {
		log.Println(err)
		respondWithError(w, http.StatusInternalServerError, "database error")
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(payload.Password))

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid password")
		return
	}

	respondWithJSON(w, http.StatusFound, POSTAPIUsersResponse{
		ID: user.ID,
		Email: user.Email,
	})
}

