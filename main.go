package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"todo-backend-go/db"
	"todo-backend-go/handler"
)

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func main() {
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// Initialize database connection
	db.InitDB()

	r := mux.NewRouter()

	// Todo API's
	r.HandleFunc("/todos", handler.GetTodos).Methods("GET")

	r.HandleFunc("/todos", handler.CreateTodo).Methods("POST")
	r.HandleFunc("/todos/{id}", handler.UpdateTodo).Methods("PATCH")
	
	r.HandleFunc("/todos/{id}", handler.DeleteTodo).Methods("DELETE")

	port := getEnv("SERVER_PORT", "8080")
	serverAddr := fmt.Sprintf(":%s", port)
	
	log.Printf("Server running on %s", serverAddr)
	log.Fatal(http.ListenAndServe(serverAddr, r))
}
