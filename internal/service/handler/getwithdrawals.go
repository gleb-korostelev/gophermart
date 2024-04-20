package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/gleb-korostelev/gophermart.git/internal/config"
	"github.com/gleb-korostelev/gophermart.git/tools/logger"
)

func (svc *APIService) GetWithdrawals(w http.ResponseWriter, r *http.Request) {
	login, ok := r.Context().Value(config.UserContextKey).(string)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	withdrawals, err := svc.store.GetWithdrawals(context.Background(), login)
	if err != nil {
		logger.Infof("Internal server error: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if len(withdrawals) == 0 {
		logger.Infof("No withdrawals found")
		w.WriteHeader(http.StatusNoContent)
		return
	}
	logger.Infof("Withdrawals were successfully read")
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(withdrawals)
}
