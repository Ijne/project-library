package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Ijne/project-library/search_app/internal/storage"
)

func FindBook(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Query().Get("q")
	fmt.Printf("title: -%s-\n", title)

	books, err := storage.GetBookByTitle(title)
	if err != nil {
		fmt.Println(err)
		return
	}

	w.Header().Set("Content-type", "application/json")
	json.NewEncoder(w).Encode(books)
}
