package storage

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/Ijne/project-library/auth_and_catalog_app/internal/models"

	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

var (
	dbPool *pgxpool.Pool
	once   sync.Once
)

func hashPassword(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	return string(hashedBytes), nil
}

func validatePassword(hashedPassword, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)) == nil
}

func initDB() {
	connStr := "postgres://go_server:go_server@localhost:5432/projectlibrary?application_name=my_go_app"

	config, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to parse config: %v\n", err)
		os.Exit(1)
	}

	config.MaxConns = 100
	config.MinConns = 5
	config.MaxConnLifetime = time.Hour
	config.MaxConnIdleTime = time.Minute * 30
	config.HealthCheckPeriod = time.Minute

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create connection pool: %v\n", err)
		os.Exit(1)
	}

	dbPool = pool
	fmt.Println("Successfully created database connection pool!")
}

func AddUser(name, email, password string) (int32, error) {
	once.Do(initDB)

	if err := dbPool.QueryRow(context.Background(), "select 1 from users where email=$1", email).Scan(new(interface{})); err == nil {
		return 0, fmt.Errorf("user with email %s already exists", email)
	}

	var id int32
	password, err := hashPassword(password)
	if err != nil {
		return 0, fmt.Errorf("failed to hash password: %w", err)
	}

	err = dbPool.QueryRow(context.Background(), "INSERT INTO users (name, email, password) VALUES ($1, $2, $3) RETURNING id", name, email, password).Scan(&id)
	if err != nil {
		fmt.Println("Failed to insert user:", err)
		return 0, fmt.Errorf("failed to insert user: %w", err)
	}

	return id, nil
}

func GetUserByEmail(email, password string) (int32, string, error) {
	once.Do(initDB)

	var id int32
	var hashedPassword, name string

	err := dbPool.QueryRow(context.Background(), "SELECT id, name, password FROM users WHERE email = $1", email).Scan(&id, &name, &hashedPassword)
	if err != nil {
		return 0, "", fmt.Errorf("failed to get user by email: %w", err)
	}

	if !validatePassword(hashedPassword, password) {
		return 0, "", fmt.Errorf("invalid email or password")
	}

	return id, name, nil
}

func GetBooks(page int) ([]models.Book, error) {
	once.Do(initDB)

	var books []models.Book

	rows, err := dbPool.Query(context.Background(), "SELECT id, title, author, status, user_id FROM books LIMIT 20 OFFSET $1", 20*(page-1))
	if err != nil {
		return nil, fmt.Errorf("failed to get books: %w", err)
	}

	for rows.Next() {
		var b models.Book
		if err := rows.Scan(&b.ID, &b.Title, &b.Author, &b.Status, &b.UserID); err != nil {
			fmt.Println("Failed to scan book row:", err)
		}
		books = append(books, b)
	}

	return books, nil
}

func GetTotalBooksCount() (int, error) {
	once.Do(initDB)

	var count int
	err := dbPool.QueryRow(context.Background(), "SELECT COUNT(*) FROM books").Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to get total books count: %w", err)
	}

	return count, nil
}
