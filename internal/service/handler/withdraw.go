package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/gleb-korostelev/gophermart.git/internal/config"
	"github.com/gleb-korostelev/gophermart.git/internal/models"
	"github.com/gleb-korostelev/gophermart.git/internal/service/utils"
	"github.com/gleb-korostelev/gophermart.git/tools/logger"
)

func (svc *APIService) Withdraw(w http.ResponseWriter, r *http.Request) {
	login, ok := r.Context().Value(config.UserContextKey).(string)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	var req models.WithdrawRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil || len(req.Order) == 0 {
		logger.Infof("Invalid request format: %v", err)
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	if !utils.IsValidOrderID(req.Order) {
		logger.Infof("Invalid order ID format")
		http.Error(w, "Invalid order ID format", http.StatusUnprocessableEntity)
		return
	}

	err = svc.store.ProcessWithdrawal(context.Background(), login, req)
	if err != nil {
		if err == config.ErrNoFunds {
			logger.Infof("Insufficient funds")
			http.Error(w, "Insufficient funds", http.StatusPaymentRequired)
			return
		}
		logger.Infof("Internal server error: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	logger.Infof("Assets were succesfully withdrawn")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Withdrawal successful"))
}
