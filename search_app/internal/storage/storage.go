package storage

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/Ijne/project-library/search_app/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	dbPool *pgxpool.Pool
	once   sync.Once
)

func initDb() {
	connStr := "postgres://search_app:search_app@localhost:5432/projectlibrary?application_name=search_app"

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
