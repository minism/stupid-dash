package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"text/template"
)

func main() {
	// Init template.
	tmpl, err := template.ParseFiles("index.html")
	if err != nil {
		panic(err)
	}

	// Get variables
	port := getEnvOr("PORT", "3000")
	hostname := getEnvOr("DASHBOARD_HOSTNAME", "http://localhost")
	title := getEnvOr("DASHBOARD_TITLE", "Dashboard")

	// Init http server.
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl.Execute(w, getTemplateData(hostname, title))
	})
	fmt.Printf("Listening on %v\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", port), nil))
}

func getEnvOr(key string, defaultValue string) string {
	value := os.Getenv(key)
	if len(value) > 0 {
		return value
	}
	return defaultValue
}
