package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthEndpoint(t *testing.T){
	t.Run("test healthz endpoint", func(t *testing.T){
		request, _ := http.NewRequest(http.MethodGet, "/api/v1/healthz", nil)
		response := httptest.NewRecorder()

		apiCfg := apiConfig{}

		apiCfg.healthz(response, request)

		got := response.Body.String()
		want := `{"status":"ok"}`

		if got != want {
			t.Errorf("got %q wanted %q", got, want)
		}
	})
}
