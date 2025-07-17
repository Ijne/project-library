package storage

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/Ijne/project-library/search_app/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

var (
	dbPool *pgxpool.Pool
	once   sync.Once
)

func initDb() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error with load .env")
	}

	DB_HOST := os.Getenv("DB_HOST")
	DB_PORT := os.Getenv("DB_PORT")
	DB_USER := os.Getenv("DB_USER")
	connStr := fmt.Sprintf("postgres://%s:search_app@%s:%s/projectlibrary?application_name=search_app", DB_USER, DB_HOST, DB_PORT)

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
	fmt.Println("Successfully created database connection pool!(search_app)")
}

func GetBookByTitle(title string) ([]models.Book, error) {
	once.Do(initDb)

	var books []models.Book
	rows, err := dbPool.Query(context.Background(), "SELECT * FROM books WHERE title LIKE $1", "%"+title+"%")
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
