package main

import (
	"app/handlers"
	"app/repository"
	"app/server"
	"app/utils"

	"go.uber.org/zap"
)

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()

	db, err := utils.Db_Init()
	if err != nil {
		logger.Fatal("database connection failed", zap.Error(err))
	}
	logger.Info("database connection success")
	repo := repository.NewRepository(db)
	handler := handlers.NewHandler(repo, logger)
	err = server.NewServer(handler)
	if err != nil {
		logger.Fatal("server failed to start", zap.Error(err))
	}

}
