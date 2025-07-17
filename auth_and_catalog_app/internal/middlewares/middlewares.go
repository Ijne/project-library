package middlewares

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/Ijne/project-library/auth_and_catalog_app/media/templates"
	"github.com/joho/godotenv"

	"github.com/golang-jwt/jwt"
)

func JWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

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

		ctx := context.WithValue(r.Context(), "username", string(username))

		next.ServeHTTP(w, r.WithContext(ctx))
	})
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
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	jwtSecret := []byte(os.Getenv("JWTsecret"))

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
