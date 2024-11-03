package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"arif/config"
	"arif/handler"
)

// Directory to save uploaded PDF files
const uploadDir = "./uploads/"

func main() {
	err := ParseConfig(".env")
	if err != nil {
		log.Fatal(err)
	}

	// Create the uploads directory if it doesn't exist
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		log.Fatalf("Could not create upload directory: %v", err)
	}

	http.HandleFunc("/upload", handler.UploadPDFHandler)

	fmt.Println("Starting server on :8098...")
	if err := http.ListenAndServe(":8098", nil); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

func ParseConfig(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	dbConfig := &config.DatabaseConfig{}
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Split the line into key and value
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue // or return an error if the format is strict
		}

		key := strings.TrimSpace(parts[0])
		value := strings.Trim(strings.TrimSpace(parts[1]), `"'`)

		switch key {
		case "DATABASE_HOST":
			dbConfig.Host = value
		case "DATABASE_PORT":
			dbConfig.Port = value
		case "DATABASE_USER":
			dbConfig.User = value
		case "DATABASE_PASSWORD":
			dbConfig.Password = value
		case "DATABASE_NAME":
			dbConfig.Name = value
		case "STORAGE_ACCESS_KEY":
			config.Cfg.StorageAccessKey = value
		case "STORAGE_SECRET_KEY":
			config.Cfg.StorageSecretKey = value
		case "CHAT_GPT_API_KEY":
			config.Cfg.ChatGPTKey = value
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	config.Cfg.DatabaseConfig = dbConfig

	return nil
}
