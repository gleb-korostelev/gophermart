package handler

import (
	"github.com/gleb-korostelev/gophermart.git/internal/service"
	"github.com/gleb-korostelev/gophermart.git/internal/storage"
	"github.com/gleb-korostelev/gophermart.git/internal/workerpool"
)

type APIService struct {
	store  storage.Storage
	worker *workerpool.DBWorkerPool
}

func NewAPIService(store storage.Storage, worker *workerpool.DBWorkerPool) service.APIServiceI {
	return &APIService{
		store:  store,
		worker: worker,
	}
}
