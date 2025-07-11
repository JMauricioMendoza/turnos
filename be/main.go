package main

import (
	"fmt"
	"log"
	"os"
	"turnos-api/database"
	"turnos-api/routes"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No se pudo cargar el archivo .env")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	database.ConectarDB()

	r := routes.SetupRouter()

	fmt.Printf("Servidor corriendo en http://localhost:%s\n", port)
	r.Run(":" + port)
}
