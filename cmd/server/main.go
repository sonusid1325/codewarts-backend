package main

import (
	"log"
	"net/http"
	"os"
	"strconv"

	"linuxmodule/pkg/container"
	"linuxmodule/pkg/db"
	"linuxmodule/pkg/router"
)

func main() {
	log.Println("Starting Linux Terminal Quest Backend...")

	// 1. Load Configurations from Env
	dbHost := getEnv("DB_HOST", "localhost")
	dbPortStr := getEnv("DB_PORT", "5432")
	dbUser := getEnv("DB_USER", "postgres")
	dbPass := getEnv("DB_PASSWORD", "postgres")
	dbName := getEnv("DB_NAME", "learn_linux")
	dbSSLMode := getEnv("DB_SSLMODE", "disable")
	port := getEnv("PORT", "8080")

	dbPort, err := strconv.Atoi(dbPortStr)
	if err != nil {
		log.Fatalf("Invalid DB_PORT value: %v", err)
	}

	// 2. Initialize Database & Run Migrations
	_, err = db.InitDB(dbHost, dbPort, dbUser, dbPass, dbName, dbSSLMode)
	if err != nil {
		log.Fatalf("Database initialization failed: %v", err)
	}

	// 3. Initialize Docker Client
	err = container.InitDocker()
	if err != nil {
		log.Fatalf("Docker initialization failed: %v", err)
	}

	// 4. Setup Router
	handler := router.SetupRouter()

	// 5. Start Server
	log.Printf("Server is starting on port %s...", port)
	err = http.ListenAndServe(":"+port, handler)
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
