package storage

import (
	"context"

	"github.com/gleb-korostelev/gophermart.git/internal/models"
)

type Storage interface {
	Ping(ctx context.Context) (int, error)
	Close() error
	Register(ctx context.Context, userCred models.User) error
	Auth(ctx context.Context, userCred models.User) error
	Orders(ctx context.Context, login, orderID string) (string, bool, error)
	GetOrders(ctx context.Context, login string) ([]models.OrdersData, error)
	GetBalances(ctx context.Context, login string) (models.BalanceData, error)
	ProcessWithdrawal(ctx context.Context, login string, req models.WithdrawRequest) error
	GetWithdrawals(ctx context.Context, login string) ([]models.Withdraws, error)
	GetOrderByNumber(ctx context.Context, orderInfo models.OrderResponse, resultChan chan<- models.OrderResponse) func(ctx context.Context) error
	UpdateOrderInfo(ctx context.Context, login string, orders models.OrderResponse) error
}
