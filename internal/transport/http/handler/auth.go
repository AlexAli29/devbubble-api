package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func NewAuthHandler() chi.Router {
	r := chi.NewRouter()

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Auth"))
	})

	return r
}
