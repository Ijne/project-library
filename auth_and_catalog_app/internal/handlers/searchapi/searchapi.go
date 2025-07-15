package searchapi

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func Search(w http.ResponseWriter, r *http.Request) {
	url := "http://localhost:8070/"

	resp, err := http.Get(url)
	if err != nil {
		http.Error(w, "Error", http.StatusBadRequest)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Error", http.StatusConflict)
		return
	}

	fmt.Printf("Статус-код: %d\n", resp.StatusCode)
	fmt.Printf("Тело ответа: %s\n", body)

}
