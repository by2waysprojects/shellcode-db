package main

import (
	"log"
	"net/http"
	"os"
	"shellcode-db/controllers"
	"shellcode-db/routes"
	"shellcode-db/services"
)

func main() {
	port := os.Getenv("SERVER_PORT")
	databaseURL := os.Getenv("NEO4J_DB")
	user := os.Getenv("NEO4J_USER")
	password := os.Getenv("NEO4J_PASSWORD")

	neo4jService := services.NewNeo4jService(databaseURL, user, password)
	defer neo4jService.Close()

	apiService := services.NewExploitDbService()
	dataLoaderService := services.NewDataLoaderService(apiService, neo4jService)
	shellCodeController := controllers.NewShellCodeController(neo4jService, dataLoaderService)

	router := routes.SetupRoutes(shellCodeController)

	log.Println("Server running on http://localhost:" + port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
