package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
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
			showLoginForm(w)
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
		w.Header().Set("Content-Type", "text/html")
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

func showLoginForm(w http.ResponseWriter) {
	html := `
		<form id="loginForm">
			<input type="text" name="username" placeholder="Username" required>
			<input type="password" name="password" placeholder="Password" required>
			<button type="button" onclick="submitForm()">Login</button>
		</form>

		<script>
			function submitForm() {
				const form = document.getElementById("loginForm");
				const data = {
					username: form.username.value,
					password: form.password.value
				};

				fetch("/login", {
					method: "POST",
					headers: {
						"Content-Type": "application/json"
					},
					body: JSON.stringify(data)
				})
				.then(response => {
					if (!response.ok) {
						return response.json().then(err => { throw err; });
					}
					return response.json();
				})
				.then(data => {
					alert(data.message);
					window.location.href = data.redirect;
				})
				.catch(error => {
					alert(error.message || "Login failed");
				});
			}
		</script>
	`
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(html))
}

// LoginHandler functions
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		var data struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}

		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}

		if data.Username != "admin" || data.Password != "password" {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}

		fmt.Println("User logged in:", data.Username)

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"username": data.Username,
			"exp":      time.Now().Add(24 * time.Hour).Unix(),
		})
		tokenString, err := token.SignedString(jwtSecret)
		if err != nil {
			http.Error(w, "Failed to create token", http.StatusInternalServerError)
			return
		}

		fmt.Println("Generated token:", tokenString)

		http.SetCookie(w, &http.Cookie{
			Name:     "token",
			Value:    tokenString,
			HttpOnly: true,
			Secure:   true,
			Path:     "/",
		})

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"message":  "Login successful",
			"redirect": "/",
		})

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
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
