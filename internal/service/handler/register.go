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

func (svc *APIService) Register(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "Invalid Content-Type", http.StatusBadRequest)
		return
	}
	var newUser models.User
	err := json.NewDecoder(r.Body).Decode(&newUser)
	if err != nil {
		logger.Infof("Bad request: %v", err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	if newUser.Login == "" || newUser.Password == "" {
		logger.Infof("Login and password required: %v", err)
		http.Error(w, "Login and password required", http.StatusBadRequest)
		return
	}
	err = svc.store.Register(context.Background(), newUser)
	if err != nil {
		if err == config.ErrLoginExists {
			logger.Infof("Login already taken: %v", err)
			http.Error(w, "Login already taken", http.StatusConflict)
			return
		}
		logger.Infof("Internal server error: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	utils.SetJWTInCookie(w, newUser.Login)
	w.WriteHeader(http.StatusOK)
	logger.Infof("User registered and authenticated")
	w.Write([]byte("User registered and authenticated"))
}
