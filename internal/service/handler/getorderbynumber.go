package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/gleb-korostelev/gophermart.git/internal/models"
	"github.com/gleb-korostelev/gophermart.git/internal/workerpool"
	"github.com/gleb-korostelev/gophermart.git/tools/logger"
	"github.com/go-chi/chi/v5"
)

func (svc *APIService) GetOrderByNumber(w http.ResponseWriter, r *http.Request) {
	number := chi.URLParam(r, "number")
	if number == "" {
		logger.Infof("This URL doesn't exist")
		http.Error(w, "This URL doesn't exist", http.StatusBadRequest)
		return
	}

	resultChan := make(chan models.OrderResponse, 1)

	svc.worker.AddTask(workerpool.Task{
		Action: func(ctx context.Context) error {
			err := svc.store.GetOrderByNumber(context.Background(), number, resultChan)
			if err != nil {
				logger.Errorf("Internal server error %v", err)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
			}
			return nil
		},
	})

	select {
	case response := <-resultChan:
		if response.Order == "" {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	case <-r.Context().Done():
		http.Error(w, "Too many requests", http.StatusTooManyRequests)
	}
}
