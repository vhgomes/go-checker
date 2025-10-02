package main

import (
	"context"
	"go-checker/internal/config"
	"go-checker/internal/server"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

		<-sigChan
		log.Println("🚦 Encerrando servidor...")
		cancel()
	}()

	db := config.InitDB()
	router := server.SetupRouter(ctx, db)

	log.Println("Server Running on :8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
