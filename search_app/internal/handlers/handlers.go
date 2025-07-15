package handlers

import "net/http"

func FindBook(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("response from search_app"))
}
