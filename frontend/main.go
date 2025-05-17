package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"
)

type DataResponse struct {
	ID        int       `json:"id"`
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
	Source    string    `json:"source"`
}

// Stringify DataResponse
func (d *DataResponse) String() string {
	return fmt.Sprintf(`{ID: '%d', Message: '%s', Timestamp: '%s', Source:'%s'}`, d.ID, d.Message, d.Timestamp, d.Source)
}

var logger *slog.Logger
var backendServiceURL string

func init() {
	logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))

	backendServiceURL = os.Getenv("BACKEND_URL")
	if backendServiceURL == "" {
		backendServiceURL = "http://localhost:8081"
		logger.Info("BACKEND_URL is not set, using default: " + backendServiceURL)
	} else {
		logger.Info("BACKEND_URL is set to " + backendServiceURL)
	}
}

// homeHandler handles '/' path.
func homeHandler(w http.ResponseWriter, r *http.Request) {
	tmplPath := filepath.Join("templates", "index.html")
	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		logger.Error(fmt.Sprintf("Could not parse template %s: %v", tmplPath, err))
		http.Error(w, "Internal Server Error: could not parse template.", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, nil)
	if err != nil {
		logger.Error(fmt.Sprintf("Could not execute template %s: %v", tmplPath, err))
		http.Error(w, "Internal Server Error: could not execute template.", http.StatusInternalServerError)
	}
}

// dataHandler handles '/data' path.
func dataHandler(w http.ResponseWriter, r *http.Request) {
	apiUrl, err1 := url.JoinPath(backendServiceURL, "api", "data")
	if err1 != nil {
		logger.Error(fmt.Sprintf("Could not join path: %v", err1))
	}

	logger.Info("Fetching data from " + apiUrl)

	response, err := http.Get(apiUrl)
	if err != nil {
		logger.Error(fmt.Sprintf("Could not fetch data from %s: %v", apiUrl, err))
	}

	// ensure the response body is closed when the function returns.
	defer response.Body.Close()

	// check if the request was successful (status code 200 OK)
	if response.StatusCode != http.StatusOK {
		logger.Error(fmt.Sprintf("Request failed with status code: %d", response.StatusCode))
		http.Error(w, "Internal Server Error: could not fetch data from API: "+apiUrl, http.StatusInternalServerError)
		return
	}

	// read the response body.
	body, err := io.ReadAll(response.Body)
	if err != nil {
		logger.Error(fmt.Sprintf("Could not read response body: %v", err))
	}

	var data DataResponse
	err = json.Unmarshal(body, &data)
	if err != nil {
		logger.Error(fmt.Sprintf("Could not unmarshal response body: %v", err))
	}

	logger.Info("Retrieve data from API: " + data.String())

	// craft the response to browser
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if err := json.NewEncoder(w).Encode(data); err != nil {
		logger.Error(fmt.Sprintf("Error encoding JSON response: %v", err))
		http.Error(w, fmt.Sprintf("Error encoding JSON response: %v", err), http.StatusInternalServerError)
	}
}

func main() {
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/data", dataHandler)

	port := "8080"
	logger.Info("Starting server on port " + port + "...")

	// start the server
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		logger.Error(fmt.Sprintf("Could not start server: %s", err))
	}
}
