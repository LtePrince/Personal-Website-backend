package main

import (
	"fmt"
	"net/http"

	"github.com/LtePrince/Personal-Website-backend/internal/handlers"
)

func main() {
	http.HandleFunc("/", handlers.Handler)
	fmt.Println("Server is listening on port 8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Printf("Error starting server: %s\n", err)
	}
}
