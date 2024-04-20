package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/gleb-korostelev/gophermart.git/internal/config"
	"github.com/gleb-korostelev/gophermart.git/tools/logger"
)

func (svc *APIService) GetOrders(w http.ResponseWriter, r *http.Request) {
	login, ok := r.Context().Value(config.UserContextKey).(string)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	orders, err := svc.store.GetOrders(context.Background(), login)
	if err != nil {
		logger.Infof("Internal server error: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if len(orders) == 0 {
		logger.Infof("Orders have no content")
		w.WriteHeader(http.StatusNoContent)
		return
	}
	logger.Infof("Orders were successfully read")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(orders)
}
