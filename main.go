package main

import (
	"IO_bound_task_service/internal/routes"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	httpSwagger "github.com/swaggo/http-swagger"
)

func main() {
	// Load .env
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	//setup port
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080" // default port
	}

	r := mux.NewRouter()

	// Routes
	r.HandleFunc("/Create_task", routes.CreateTask).Methods("POST")
	r.HandleFunc("/task_status/{task_id}", routes.GetTaskStatus).Methods("GET")

	// Swagger
	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	log.Printf("Server started on port :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
