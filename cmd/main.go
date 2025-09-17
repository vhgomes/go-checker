package main

import (
	"go-checker/internal/config"
	"go-checker/internal/server"
	"log"
)

func main() {
	db := config.InitDB()
	router := server.SetupRouter(db)

	log.Println("Server Running on :8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
