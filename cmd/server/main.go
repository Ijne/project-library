package main

import (
	"net/http"
	"projectlibrary/internal/handlers"
	"projectlibrary/internal/middlewares"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
)

func main() {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.With(middlewares.JWTMiddleware).Handle("/", http.HandlerFunc(handlers.HomeHandler))
	r.Handle("/register", http.HandlerFunc(handlers.RegisterHandler))
	r.Handle("/login", http.HandlerFunc(handlers.LoginHandler))
	r.Handle("/logout", http.HandlerFunc(handlers.LogoutHandler))

	err := http.ListenAndServe(":8080", r)
	if err != nil {
		panic(err)
	}
}
