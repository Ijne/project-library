package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/Ijne/project-library/auth_and_catalog_app/internal/handlers/auth"
	"github.com/Ijne/project-library/auth_and_catalog_app/internal/handlers/catalog"
	"github.com/Ijne/project-library/auth_and_catalog_app/internal/handlers/searchapi"
	"github.com/Ijne/project-library/auth_and_catalog_app/internal/middlewares"
	"github.com/joho/godotenv"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("ERROR")
		return
	}

	PORT := os.Getenv("PORT")

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.With(middlewares.JWTMiddleware).Handle("/", http.HandlerFunc(auth.HomeHandler))
	r.With(middlewares.JWTMiddleware).Handle("/catalog", http.HandlerFunc(catalog.CatalogHandler))
	r.Handle("/register", http.HandlerFunc(auth.RegisterHandler))
	r.Handle("/login", http.HandlerFunc(auth.LoginHandler))
	r.Handle("/logout", http.HandlerFunc(auth.LogoutHandler))
	r.Handle("/api/search-books", http.HandlerFunc(searchapi.Search))

	r.Handle("/media/static/*", http.StripPrefix("/media/static/", http.FileServer(http.Dir("media/static/"))))
	err = http.ListenAndServe(fmt.Sprintf(":%s", PORT), r)
	if err != nil {
		panic(err)
	}
}
