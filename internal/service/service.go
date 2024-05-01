package service

import (
	"net/http"

	"github.com/gleb-korostelev/gophermart.git/internal/models"
)

type APIServiceI interface {
	Register(w http.ResponseWriter, r *http.Request)
	Login(w http.ResponseWriter, r *http.Request)
	Orders(w http.ResponseWriter, r *http.Request)
	GetOrders(w http.ResponseWriter, r *http.Request)
	GetBalance(w http.ResponseWriter, r *http.Request)
	Withdraw(w http.ResponseWriter, r *http.Request)
	GetWithdrawals(w http.ResponseWriter, r *http.Request)
	CheckOrderStatus(login string, orders []models.OrdersData) []models.OrdersData
}
