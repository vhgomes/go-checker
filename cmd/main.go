package main

import (
	"context"
	"go-checker/internal/config"
	"go-checker/internal/server"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"
)

func main() {
	config.InitLogger()
	defer zap.L().Sync()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

		<-sigChan
		zap.L().Info("🚦 Encerrando servidor...")
		cancel()
	}()

	db := config.InitDB()
	redis := config.InitRedis()

	router := server.SetupRouter(ctx, db, redis)

	port := config.GetEnvOrDefault("PORT", "8080")
	zap.L().Info("Server Running on :" + port)
	if err := router.Run(":" + port); err != nil {
		zap.L().Fatal("Erro ao iniciar o servidor", zap.Error(err))
	}
}
