package main

import (
	"net/http"
	"projectlibrary/internal/handlers"

	"github.com/go-chi/chi/v5"
)

func main() {
	r := chi.NewRouter()

	r.Handle("/", http.HandlerFunc(handlers.HomeHandler))
	r.Handle("/login", http.HandlerFunc(handlers.LoginHandler))
	r.Handle("/logout", http.HandlerFunc(handlers.LogoutHandler))

	err := http.ListenAndServe(":8080", r)
	if err != nil {
		panic(err)
	}
}
