package catalog

import (
	"html/template"
	"net/http"
	"strconv"

	"github.com/Ijne/project-library/auth_and_catalog_app/internal/models"
	"github.com/Ijne/project-library/auth_and_catalog_app/internal/storage"
)

var (
	TotalBooksCount int
)

func CatalogHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		page, err := strconv.Atoi(r.URL.Query().Get("page"))
		if err != nil {
			http.Error(w, "Invalid page number", http.StatusBadRequest)
			return
		}
		TotalBooksCount, _ = storage.GetTotalBooksCount()

		books, err := storage.GetBooks(page)
		if err != nil {
			http.Error(w, "Failed to retrieve books", http.StatusInternalServerError)
			return
		}

		data := struct {
			Books       []models.Book
			CurrentPage int
			TotalPages  int
			PrevPage    int
			NextPage    int
		}{
			Books:       books,
			CurrentPage: page,
			TotalPages:  (TotalBooksCount + 19) / 20,
			PrevPage:    page - 1,
			NextPage:    page + 1,
		}

		renderTemplate(w, "catalog.html", data)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
}

func renderTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
	templates := template.Must(template.ParseFiles(
		"media/templates/base.html",
		"media/templates/"+tmpl,
	))

	err := templates.ExecuteTemplate(w, "base.html", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
