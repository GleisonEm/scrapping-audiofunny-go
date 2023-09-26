package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gleisonem/scrapping-audiofunny-go/controllers"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Erro ao carregar o arquivo .env")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8087"
	}

	r := mux.NewRouter()
	r.HandleFunc("/search", controllers.SearchHandler).Methods("GET")

	http.Handle("/", r)

	address := ":" + port
	fmt.Printf("Servidor escutando na porta %s...\n", port)
	http.ListenAndServe(address, nil)
}
