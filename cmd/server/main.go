package main

import (
	"net/http"
	"projectlibrary/internal/handlers/auth"
	"projectlibrary/internal/handlers/catalog"
	"projectlibrary/internal/middlewares"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
)

func main() {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.With(middlewares.JWTMiddleware).Handle("/", http.HandlerFunc(auth.HomeHandler))
	r.With(middlewares.JWTMiddleware).Handle("/catalog", http.HandlerFunc(catalog.CatalogHandler))
	r.Handle("/register", http.HandlerFunc(auth.RegisterHandler))
	r.Handle("/login", http.HandlerFunc(auth.LoginHandler))
	r.Handle("/logout", http.HandlerFunc(auth.LogoutHandler))

	r.Handle("/media/static/*", http.StripPrefix("/media/static/", http.FileServer(http.Dir("media/static/"))))
	err := http.ListenAndServe(":8080", r)
	if err != nil {
		panic(err)
	}
}
