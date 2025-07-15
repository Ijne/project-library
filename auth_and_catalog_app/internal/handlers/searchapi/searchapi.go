package searchapi

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
)

func Search(w http.ResponseWriter, r *http.Request) {
	title := url.QueryEscape(r.URL.Query().Get("q"))
	fmt.Println("title (from auth_and_catalog_app):", title)

	url := fmt.Sprintf("http://localhost:8070/?q=%s", title)

	resp, err := http.Get(url)
	if err != nil {
		http.Error(w, "Error", http.StatusBadRequest)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Error", http.StatusConflict)
		return
	}

	fmt.Printf("Статус-код: %d\n", resp.StatusCode)
	fmt.Printf("Тело ответа: %s\n", body)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.StatusCode)
	w.Write(body)

}
