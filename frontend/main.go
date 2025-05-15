package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

type TemplateData struct {
	BackendURL string
}

func main() {
	port := "8080"

	backendServiceURL := os.Getenv("BACKEND_URL")
	if backendServiceURL == "" {
		backendServiceURL = "http://localhost:8081"
		log.Printf("BACKEND_URL env variable is not set, using default: %s", backendServiceURL)
	} else {
		log.Printf("Using BACKEND_URL from environment: %s", backendServiceURL)
	}

	// handler for the root path.
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmplPath := filepath.Join("templates", "index.html")

		// parse the HTML template
		tmpl, err := template.ParseFiles(tmplPath)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error parsing template: %v", err), http.StatusInternalServerError)
			log.Printf("Error parsing template %s: %v", tmplPath, err)
			return
		}

		// Data to pass to the template
		data := TemplateData{
			BackendURL: backendServiceURL,
		}

		// Execute the template, writing the output to the response writer.
		err = tmpl.Execute(w, data)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error executing template: %v", err), http.StatusInternalServerError)
			log.Printf("Error executing template %s: %v", tmplPath, err)
		}
	})

	log.Printf("Frontend service starting on port %s", port)
	log.Printf("Serving UI that will try to connect to backend at: %s", backendServiceURL)

	// start the HTTP server
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Failed to start frontend server: %v", err)
	}
}
