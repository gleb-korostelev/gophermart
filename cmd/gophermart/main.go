package main

import (
	"net/http"

	"github.com/gleb-korostelev/gophermart.git/internal/config"
	"github.com/gleb-korostelev/gophermart.git/internal/db/dbimpl"
	"github.com/gleb-korostelev/gophermart.git/internal/service/handler"
	"github.com/gleb-korostelev/gophermart.git/internal/service/router"
	"github.com/gleb-korostelev/gophermart.git/internal/storage/repository"
	"github.com/gleb-korostelev/gophermart.git/internal/workerpool"
	"github.com/gleb-korostelev/gophermart.git/tools/logger"
	"go.uber.org/zap"
)

func main() {
	log, _ := zap.NewProduction()
	err := config.ConfigInit()
	if err != nil {
		logger.Infof("Error in config: %v", err)
		return
	}

	database, err := dbimpl.InitDB()
	if err != nil {
		logger.Infof("Error database initialize: %v", err)
		return
	}

	store := repository.NewDBStorage(database)
	defer store.Close()
	workerPool := workerpool.NewDBWorkerPool(config.MaxRoutine)
	defer workerPool.Shutdown()
	svc := handler.NewAPIService(store, workerPool)
	r := router.RouterInit(svc, log)

	logger.Infof("Server is listening on: %s", config.ServerAddr)
	if err := http.ListenAndServe(config.ServerAddr, r); err != nil {
		logger.Fatal("Error starting server: %v", err)
	}
}
