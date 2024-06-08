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
	"strconv"
	"time"

	"url-short/internal/database"
	"url-short/internal/shortener"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type apiConfig struct {
	DB        *database.Queries
	JWTSecret string
}

type HealthResponse struct {
	Status string `json:"status"`
}

type LongURLRequest struct {
	LongURL string `json:"long_url"`
}

type LongURLResponse struct {
	ShortURL string `json:"short_url"`
}

type APIUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"Password"`
}

type APIUsersResponse struct {
	ID    int32  `json:"id"`
	Email string `json:"email"`
	Token string `json:"token"`
}

type APIUserResponseNoToken struct {
	ID    int32  `json:"id"`
	Email string `json:"email"`
}

func (apiCfg *apiConfig) healthz(w http.ResponseWriter, r *http.Request) {
	payload := HealthResponse{
		Status: "ok",
	}

	respondWithJSON(w, http.StatusOK, payload)
}

func (apiCfg *apiConfig) postLongURL(w http.ResponseWriter, r *http.Request) {
	payload := LongURLRequest{}

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

	respondWithJSON(w, http.StatusCreated, LongURLResponse{
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
	payload := APIUserRequest{}

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
		Email:     payload.Email,
		Password:  string(passwordHash),
		CreatedAt: now,
		UpdatedAt: now,
	})

	respondWithJSON(w, http.StatusCreated, APIUserResponseNoToken{
		ID:    user.ID,
		Email: user.Email,
	})
}

func (apiCfg *apiConfig) postAPILogin(w http.ResponseWriter, r *http.Request) {
	payload := APIUserRequest{}

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

	registeredClaims := jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		Issuer:    "url-short-auth",
		Subject:   strconv.Itoa(int(user.ID)),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, registeredClaims)

	signedToken, err := token.SignedString([]byte(apiCfg.JWTSecret))

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "can not create JWT")
		return
	}

	respondWithJSON(w, http.StatusFound, APIUsersResponse{
		ID:    user.ID,
		Email: user.Email,
		Token: signedToken,
	})
}

func (apiCfg *apiConfig) putAPIUsers(w http.ResponseWriter, r *http.Request, user database.User) {
	payload := APIUserRequest{}

	err := json.NewDecoder(r.Body).Decode(&payload)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "incorrect parameters for user update request")
		return
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)

	if err != nil {
		log.Println(err)
		respondWithError(w, http.StatusBadRequest, "bad password supplied from client")
		return
	}

	err = apiCfg.DB.UpdateUser(r.Context(), database.UpdateUserParams{
		Email:    payload.Email,
		Password: string(passwordHash),
		ID:       user.ID,
		UpdatedAt: time.Now(),
	})

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not update user in database")
	}

	respondWithJSON(w, http.StatusOK, APIUserResponseNoToken{
		Email: payload.Email,
		ID:    user.ID,
	})
}
