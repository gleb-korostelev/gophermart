package handler

import (
	"context"
	"io"
	"net/http"

	"github.com/gleb-korostelev/gophermart.git/internal/config"
	"github.com/gleb-korostelev/gophermart.git/internal/service/utils"
	"github.com/gleb-korostelev/gophermart.git/tools/logger"
)

func (svc *APIService) Orders(w http.ResponseWriter, r *http.Request) {
	login, ok := r.Context().Value(config.UserContextKey).(string)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	orderID, err := io.ReadAll(r.Body)
	if err != nil || len(orderID) == 0 {
		logger.Infof("Invalid request format: %v", err)
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	if !utils.IsValidOrderID(string(orderID)) {
		logger.Infof("Invalid order ID format")
		http.Error(w, "Invalid order ID format", http.StatusUnprocessableEntity)
		return
	}

	orderLogin, exists, err := svc.store.Orders(context.Background(), login, string(orderID))
	if err != nil {
		logger.Infof("Internal Error in Orders function: %v", err)
		http.Error(w, "Internal Error in Orders function", http.StatusInternalServerError)
		return
	}
	if exists {
		if orderLogin == login {
			logger.Infof("Order number already uploaded by this user: ", orderLogin)
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Order number already uploaded by this user"))
			return
		} else {
			logger.Infof("Order number already uploaded by another user: ", orderLogin)
			http.Error(w, "Order number already uploaded by another user", http.StatusConflict)
			return
		}
	}
	logger.Infof("New order number accepted for processing: ", string(orderID))
	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte("New order number accepted for processing"))
}
