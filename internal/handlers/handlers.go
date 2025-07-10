package handlers

import (
	"encoding/json"
	"net/http"
	"projectlibrary/internal/storage"
	"time"

	"github.com/golang-jwt/jwt"
)

var jwtSecret = []byte("sPfasW$#@D32as+*qwrg32da")

// HomeHandler functions
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write([]byte("Welcome, " + r.Context().Value("username").(string) + "! <a href='/logout'>Logout</a>"))
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
}

func makeCookieAfterLogin(w http.ResponseWriter, id int32, username string) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":       id,
		"username": username,
		"exp":      time.Now().Add(24 * time.Hour).Unix(),
	})
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		http.Error(w, "Failed to create token", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    tokenString,
		HttpOnly: true,
		Secure:   true,
		Path:     "/",
	})

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(map[string]string{
		"message":  "Register successful",
		"redirect": "/",
	})
}

// RegisterHandler functions
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		var data struct {
			Username string `json:"username"`
			Password string `json:"password"`
			Email    string `json:"email"`
		}

		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}

		id, err := storage.AddUser(data.Username, data.Email, data.Password)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		makeCookieAfterLogin(w, id, data.Username)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
}

// LoginHandler functions
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		var data struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}

		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}

		id, name, err := storage.GetUserByEmail(data.Email, data.Password)
		if err != nil {
			http.Error(w, "Invalid email or password", http.StatusUnauthorized)
			return
		}

		makeCookieAfterLogin(w, id, name)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}

}

// LogoutHandler functions
func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		http.SetCookie(w, &http.Cookie{
			Name:     "token",
			Value:    "",
			HttpOnly: true,
			Secure:   true,
			Path:     "/",
			MaxAge:   -1,
		})
		http.Redirect(w, r, "/", http.StatusFound)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
