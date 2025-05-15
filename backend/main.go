package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/lib/pq" // PostgreSQL driver
)

// App holds application-wide dependencies, like a database connection.
type App struct {
	DB *sql.DB
}

// DataResponse defines the structure of our API response.
type DataResponse struct {
	ID        int       `json:"id"`
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
	Source    string    `json:"source"`
}

// HealthCheckResponse defines the structure for healthcheck
type HealthCheckResponse struct {
	Status    string `json:"status"`
	DBStatus  string `json:"db_status"`
	Timestamp string `json:"timestamp"`
}

// Creates a simple 'items' table if it doesn't exist.
func (a *App) createItemsTable() error {
	if a.DB == nil {
		return fmt.Errorf("database connection is not initialized")
	}
	query := `
		CREATE TABLE IF NOT EXISTS items (
			id SERIAL PRIMARY KEY,
			name VARCHAR(100) NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);`
	_, err := a.DB.Exec(query)
	if err != nil {
		return fmt.Errorf("error creating items table: %w", err)
	}
	log.Println("Successfully ensured 'items' table exists.")

	return nil
}

// Handles requests to the /api/data endpoint
func (a *App) dataHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Received request to /api/data")

	responseData := DataResponse{
		ID:        1,
		Message:   "Default mock message",
		Timestamp: time.Now(),
		Source:    "mock",
	}

	if a.DB != nil {
		// attempt to insert a new item and retrieve its name, or just query an existing one.
		// for simplicity, let's insert a new item each time and return its details.
		var itemName string
		var itemID int
		var itemCreatedAt time.Time

		itemNameInput := fmt.Sprintf("Sample Item %d", time.Now().UnixNano())

		err := a.DB.QueryRow(
			"INSERT INTO items (name) VALUES ($1) RETURNING id, name, created_at",
			itemNameInput,
		).Scan(&itemID, &itemName, &itemCreatedAt)

		if err != nil {
			log.Printf("Error interacting with database: %v", err)
			responseData.Message = fmt.Sprintf("Error interacting with database: %v", err)
			responseData.Source = "error"
		} else {
			responseData.ID = itemID
			responseData.Message = fmt.Sprintf("Item from DB: %s", itemName)
			responseData.Timestamp = itemCreatedAt
			responseData.Source = "database"
			log.Printf("Successfully retrieved/inserted item from DB: %d, Name=%s", itemID, itemName)
		}

	} else {
		responseData.Message = "Hello from the Go Backend API (DB not connected)"
		log.Println("Database not connected, returning mock data.")
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if err := json.NewEncoder(w).Encode(responseData); err != nil {
		log.Printf("Error encoding JSON response: %v", err)
		http.Error(w, fmt.Sprintf("Error encoding JSON response: %v", err), http.StatusInternalServerError)
	}
}

// Handles requests to the /health endpoint.
func (a *App) healthHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Received request to /health")
	dbStatus := "DOWN"

	if a.DB != nil {
		if err := a.DB.Ping(); err == nil {
			dbStatus = "UP"
		} else {
			log.Printf("DB Ping failed: %v", err)
		}
	}

	response := HealthCheckResponse{
		Status:    "OK",
		DBStatus:  dbStatus,
		Timestamp: time.Now().Format(time.RFC3339),
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(response)
}

func main() {
	port := "8081" // port for the Go application
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	var app App
	var err error

	if dbHost != "" {
		connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			dbHost, dbPort, dbUser, dbPassword, dbName)

		log.Printf("Attempting to connect to database with connStr: host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			dbHost, port, dbUser, dbPassword, dbName)

		app.DB, err = sql.Open("postgres", connStr)
		if err != nil {
			log.Printf("Warning: Failed to open database connection: %v. Running without DB", err)
			app.DB = nil
		} else {
			err = app.DB.Ping()
			if err != nil {
				log.Printf("Warning: Failed to ping database: %v. Running without DB", err)
				app.DB.Close()
				app.DB = nil
			} else {
				log.Println("Successfully connected to the database")
				// create table if it doesn't exist.
				if err := app.createItemsTable(); err != nil {
					log.Fatalf("Failed to create items table: %v", err)
				}

			}
		}
	} else {
		log.Println("DB_HOST environment variable not set. Running backend service in mock data mode.")
	}

	if app.DB != nil {
		defer func() {
			log.Println("Closing database connection...")
			if err := app.DB.Close(); err != nil {
				log.Printf("Error closing database connection: %v", err)
			}
		}()
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/api/data", app.dataHandler)
	mux.HandleFunc("/health", app.healthHandler)

	log.Printf("Go API Service starting on port %s", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
