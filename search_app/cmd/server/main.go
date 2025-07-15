package main

import (
	"net/http"

	"github.com/Ijne/project-library/search_app/internal/handlers"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func main() {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Handle("/", http.HandlerFunc(handlers.FindBook))

	if err := http.ListenAndServe(":8070", r); err != nil {
		panic(err)
	}
}
