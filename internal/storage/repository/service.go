package repository

import (
	"context"
	"net/http"

	"github.com/gleb-korostelev/gophermart.git/internal/config"
	"github.com/gleb-korostelev/gophermart.git/internal/db"
	"github.com/gleb-korostelev/gophermart.git/internal/db/dbimpl"
	"github.com/gleb-korostelev/gophermart.git/internal/models"
	"github.com/gleb-korostelev/gophermart.git/internal/storage"
	"github.com/gleb-korostelev/gophermart.git/tools/logger"
	"golang.org/x/crypto/bcrypt"
)

type service struct {
	data db.DB
}

func NewDBStorage(data db.DB) storage.Storage {
	return &service{
		data: data,
	}
}

func (s *service) Ping(ctx context.Context) (int, error) {
	err := s.data.Ping(context.Background())
	if err != nil {
		logger.Errorf("Failed to connect to the database %v", err)
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil
}

func (s *service) Close() error {
	err := s.data.Close()
	if err != nil {
		return err
	}
	return nil
}

func (s *service) Register(ctx context.Context, userCred models.User) error {

	EncryptedPassword, err := bcrypt.GenerateFromPassword([]byte(userCred.Password), bcrypt.DefaultCost)
	if err != nil {
		logger.Infof("Failed to encrypt user's password: %v", err)
		return err
	}

	err = dbimpl.SaveUser(s.data, ctx, userCred.Login, string(EncryptedPassword))
	if err != nil {
		logger.Infof("Failed to save user to database: %v", err)
		return err
	}
	return nil
}

func (s *service) Auth(ctx context.Context, userCred models.User) error {
	EncryptedPassword, err := dbimpl.GetUserCred(s.data, ctx, userCred.Login)
	if err != nil {
		logger.Infof("Failed to find user in database: %v", err)
		return err
	}

	err = bcrypt.CompareHashAndPassword([]byte(EncryptedPassword), []byte(userCred.Password))
	if err != nil {
		logger.Infof("Failed to authorize: %v", config.ErrWrongPassword)
		return config.ErrWrongPassword
	}
	return nil
}

func (s *service) Orders(ctx context.Context, login, orderID string) (string, bool, error) {
	orderLogin, exists, err := dbimpl.SaveOrders(s.data, ctx, login, orderID)
	if err != nil {
		return "", false, err
	}
	return orderLogin, exists, err
}

func (s *service) GetOrders(ctx context.Context, login string) ([]models.OrdersData, error) {
	orders, err := dbimpl.GetOrders(s.data, ctx, login)
	if err != nil {
		return orders, err
	}
	return orders, nil
}

func (s *service) GetBalances(ctx context.Context, login string) (models.BalanceData, error) {
	balance, err := dbimpl.Balance(s.data, ctx, login)
	if err != nil {
		return models.BalanceData{}, err
	}
	return balance, nil
}

func (s *service) ProcessWithdrawal(ctx context.Context, login string, req models.WithdrawRequest) error {
	return dbimpl.Withdraw(s.data, ctx, login, req)
}

func (s *service) GetWithdrawals(ctx context.Context, login string) ([]models.Withdraws, error) {
	return dbimpl.GetWithdrawals(s.data, ctx, login)
}

func (s *service) GetOrderByNumber(ctx context.Context, orderInfo models.OrderResponse, resultChan chan<- models.OrderResponse) func(ctx context.Context) error {
	return dbimpl.GetOrderByNumber(s.data, ctx, orderInfo, resultChan)
}

func (s *service) UpdateOrderInfo(ctx context.Context, login string, orders models.OrderResponse) error {
	return dbimpl.UpdateOrderInfo(s.data, ctx, login, orders)
}
