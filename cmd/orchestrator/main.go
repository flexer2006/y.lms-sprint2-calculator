package main

import (
	"net/http"

	"github.com/flexer2006/y.lms-sprint2-calculator/internal/logger"

	"github.com/gorilla/mux"
)

func main() {
	logger.Info("Starting Orchestrator...")

	r := mux.NewRouter()

	r.HandleFunc("/api/v1/calculate", calculateHandler).Methods("POST")
	r.HandleFunc("/api/v1/expressions", listExpressionsHandler).Methods("GET")
	r.HandleFunc("/api/v1/expressions/{id}", getExpressionHandler).Methods("GET")
	r.HandleFunc("/internal/task", getTaskHandler).Methods("GET")
	r.HandleFunc("/internal/task", postTaskResultHandler).Methods("POST")

	http.Handle("/", r)
	logger.Info("Orchestrator is running on port 8080")
	http.ListenAndServe(":8080", nil)
}

func calculateHandler(w http.ResponseWriter, r *http.Request) {
	// Handler logic for calculating expressions
}

func listExpressionsHandler(w http.ResponseWriter, r *http.Request) {
	// Handler logic for listing expressions
}

func getExpressionHandler(w http.ResponseWriter, r *http.Request) {
	// Handler logic for getting a specific expression
}

func getTaskHandler(w http.ResponseWriter, r *http.Request) {
	// Handler logic for getting a task
}

func postTaskResultHandler(w http.ResponseWriter, r *http.Request) {
	// Handler logic for posting task results
}
