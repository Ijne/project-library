package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"projectlibrary/internal/storage"
	"projectlibrary/media/templates"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// HomeHandler functions
var jwtSecret = []byte("sPfasW$#@D32as+*qwrg32da")

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		token := extractToken(r)

		if token == "" {
			action := r.URL.Query().Get("action")
			switch action {
			case "register":
				showForm(w, templates.RegisterTemplate)
			case "login":
				showForm(w, templates.LoginTemplate)
			default:
				showForm(w, templates.RegisterTemplate)
			}
			return
		}

		claims, err := validateToken(token)
		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		username, ok := claims.(jwt.MapClaims)["username"].(string)
		if !ok || username == "" {
			http.Error(w, "Invalid token claims", http.StatusUnauthorized)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write([]byte("Welcome, " + username + "! <a href='/logout'>Logout</a>"))
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
}

func extractToken(r *http.Request) string {
	authHeader := r.Header.Get("Authorization")
	if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
		return strings.TrimPrefix(authHeader, "Bearer ")
	}

	token, err := r.Cookie("token")
	if err != nil {
		return ""
	}

	return token.Value
}

func validateToken(tokenString string) (jwt.Claims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected algorithm: %v", token.Header["alg"])
		}
		return jwtSecret, nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	return token.Claims, nil
}

func showForm(w http.ResponseWriter, html string) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(html))
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
