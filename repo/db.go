package repo

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"

	"arif/config"
)

const (
	statusCreated = "created"
)

func CreateRequest(ctx context.Context, hash, url string) error {
	// Connection string: update with your actual database credentials
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", config.Cfg.DatabaseConfig.Host, config.Cfg.DatabaseConfig.Port, config.Cfg.DatabaseConfig.User, config.Cfg.DatabaseConfig.Password, config.Cfg.DatabaseConfig.Name)

	// Open the database connection
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return fmt.Errorf("failed to open a DB connection: %v", err)
	}
	defer db.Close()

	// Verify the connection is valid
	if err = db.Ping(); err != nil {
		return fmt.Errorf("failed to ping DB: %v", err)
	}

	// SQL statement to insert a row into 'your_table' with columns 'column1' and 'column2'
	query := "INSERT INTO requests (hash, status, url) VALUES ($1, $2, $3);"

	// Execute the statement with placeholder values
	_, err = db.ExecContext(ctx, query, hash, statusCreated, url)
	if err != nil {
		return fmt.Errorf("failed to execute query: %v", err)
	}

	fmt.Println("Row inserted successfully.")
	return nil
}

func UpdateRequestStatus(ctx context.Context, hash, status string) error {
	// Connection string: update with your actual database credentials
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", config.Cfg.DatabaseConfig.Host, config.Cfg.DatabaseConfig.Port, config.Cfg.DatabaseConfig.User, config.Cfg.DatabaseConfig.Password, config.Cfg.DatabaseConfig.Name)

	// Open the database connection
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return fmt.Errorf("failed to open a DB connection: %v", err)
	}
	defer db.Close()

	// Verify the connection is valid
	if err = db.Ping(); err != nil {
		return fmt.Errorf("failed to ping DB: %v", err)
	}

	// SQL statement to insert a row into 'your_table' with columns 'column1' and 'column2'
	query := "UPDATE requests SET status = $1, updated_date = now() WHERE hash = $2;"

	// Execute the statement with placeholder values
	_, err = db.ExecContext(ctx, query, status, hash)
	if err != nil {
		return fmt.Errorf("failed to execute query: %v", err)
	}

	fmt.Println("Row inserted successfully.")
	return nil
}

func CreateEntry(ctx context.Context, hash string, urls map[int]string) error {
	// Connection string: update with your actual database credentials
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", config.Cfg.DatabaseConfig.Host, config.Cfg.DatabaseConfig.Port, config.Cfg.DatabaseConfig.User, config.Cfg.DatabaseConfig.Password, config.Cfg.DatabaseConfig.Name)

	// Open the database connection
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return fmt.Errorf("failed to open a DB connection: %v", err)
	}
	defer db.Close()

	// Verify the connection is valid
	if err = db.Ping(); err != nil {
		return fmt.Errorf("failed to ping DB: %v", err)
	}

	// SQL statement to insert a row into 'your_table' with columns 'column1' and 'column2'
	query := "INSERT INTO entries (request_hash, page, url) VALUES ($1, $2, $3);"

	// Execute the statement with placeholder values
	for page, url := range urls {
		_, err = db.ExecContext(ctx, query, hash, page, url)
		if err != nil {
			return fmt.Errorf("failed to execute query: %v", err)
		}
	}

	fmt.Println("Row inserted successfully.")
	return nil
}

func InsertExtracted(ctx context.Context, hash string, pageExtractedMap map[int]string) error {
	// Connection string: update with your actual database credentials
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", config.Cfg.DatabaseConfig.Host, config.Cfg.DatabaseConfig.Port, config.Cfg.DatabaseConfig.User, config.Cfg.DatabaseConfig.Password, config.Cfg.DatabaseConfig.Name)

	// Open the database connection
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return fmt.Errorf("failed to open a DB connection: %v", err)
	}
	defer db.Close()

	// Verify the connection is valid
	if err = db.Ping(); err != nil {
		return fmt.Errorf("failed to ping DB: %v", err)
	}

	// SQL statement to insert a row into 'your_table' with columns 'column1' and 'column2'
	query := "UPDATE entries SET extracted = $1, updated_date = now() WHERE request_hash = $2 and page = $3;"

	// Execute the statement with placeholder values
	for page, extracted := range pageExtractedMap {
		_, err = db.ExecContext(ctx, query, extracted, hash, page)
		if err != nil {
			return fmt.Errorf("failed to execute query: %v", err)
		}
	}

	fmt.Println("Row inserted successfully.")
	return nil
}

func InsertTranslated(ctx context.Context, hash string, pageTranslatedMap map[int]string) error {
	// Connection string: update with your actual database credentials
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", config.Cfg.DatabaseConfig.Host, config.Cfg.DatabaseConfig.Port, config.Cfg.DatabaseConfig.User, config.Cfg.DatabaseConfig.Password, config.Cfg.DatabaseConfig.Name)

	// Open the database connection
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return fmt.Errorf("failed to open a DB connection: %v", err)
	}
	defer db.Close()

	// Verify the connection is valid
	if err = db.Ping(); err != nil {
		return fmt.Errorf("failed to ping DB: %v", err)
	}

	// SQL statement to insert a row into 'your_table' with columns 'column1' and 'column2'
	query := "UPDATE entries SET translated = $1, updated_date = now() WHERE request_hash = $2 and page = $3;"

	// Execute the statement with placeholder values
	for page, translated := range pageTranslatedMap {
		_, err = db.ExecContext(ctx, query, translated, hash, page)
		if err != nil {
			return fmt.Errorf("failed to execute query: %v", err)
		}
	}

	fmt.Println("Row inserted successfully.")
	return nil
}
