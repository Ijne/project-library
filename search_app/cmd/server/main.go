package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Ijne/project-library/search_app/internal/handlers"
	"github.com/joho/godotenv"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error with load .env")
	}

	PORT := os.Getenv("PORT")

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Handle("/", http.HandlerFunc(handlers.FindBook))

	if err := http.ListenAndServe(fmt.Sprintf(":%s", PORT), r); err != nil {
		panic(err)
	}
}
