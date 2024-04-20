package dbimpl

import (
	"context"
	"time"

	"github.com/gleb-korostelev/gophermart.git/internal/config"
	"github.com/gleb-korostelev/gophermart.git/internal/db"
	"github.com/gleb-korostelev/gophermart.git/internal/models"
	"github.com/jackc/pgx/v5"
)

func GetUserCred(db db.DB, ctx context.Context, login string) (string, error) {
	var password string
	var isDeleted bool
	sql := `SELECT password FROM user_data WHERE login = $1`
	err := db.QueryRow(ctx, sql, login).Scan(&password)
	if err != nil {
		return "", err
	}
	if isDeleted {
		return "", config.ErrGone
	}
	return password, nil
}

func GetOrdersData(db db.DB, ctx context.Context, login string) ([]models.OrdersData, error) {
	sql := `
	SELECT order_id, status, accrual, uploaded_at
	FROM orders
	WHERE login=$1
	ORDER BY uploaded_at ASC
	`
	rows, err := db.Query(ctx, sql, login)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []models.OrdersData
	for rows.Next() {
		var order models.OrdersData
		var accrual *float64
		var uploadedAt time.Time
		if err := rows.Scan(&order.Number, &order.Status, &accrual, &uploadedAt); err != nil {
			return nil, err
		}
		if accrual != nil {
			order.Accrual = *accrual
		}
		order.UploadedAt = uploadedAt.Format(time.RFC3339)
		orders = append(orders, order)
	}
	return orders, nil
}

func Balance(db db.DB, ctx context.Context, login string) (models.BalanceData, error) {
	sql := `
	SELECT current_balance, withdrawn
	FROM balances
	WHERE login=$1
	`
	var balance models.BalanceData
	err := db.QueryRow(ctx, sql, login).Scan(&balance.Current, &balance.Withdrawn)
	if err != nil {
		return models.BalanceData{}, err
	}
	return balance, nil
}

func GetWithdrawals(db db.DB, ctx context.Context, login string) ([]models.Withdraws, error) {
	sql := `
	SELECT order_id, sum, processed_at
	FROM withdrawals
	WHERE login=$1
	`
	rows, err := db.Query(ctx, sql, login)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var withdraws []models.Withdraws
	for rows.Next() {
		var withdraw models.Withdraws
		var processedAt time.Time
		if err := rows.Scan(&withdraw.Order, &withdraw.Sum, &processedAt); err != nil {
			return nil, err
		}
		withdraw.ProcessedAt = processedAt.Format(time.RFC3339)
		withdraws = append(withdraws, withdraw)
	}
	return withdraws, nil
}

func GetOrder(db db.DB, ctx context.Context, order string, resultChan chan<- models.OrderResponse) func(ctx context.Context) error {
	return func(ctx context.Context) error {
		var response models.OrderResponse
		var accrual *float64
		sql := `
	SELECT order_id, status, accrual
	FROM orders
	WHERE order_id = $1
	`
		err := db.QueryRow(ctx, sql, order).Scan(&response.Order, &response.Status, &accrual)
		if err == pgx.ErrNoRows {
			return config.ErrNotFound
		} else if err != nil {
			return err
		}

		if accrual != nil {
			response.Accrual = *accrual
		}
		select {
		case resultChan <- response:
		case <-ctx.Done():
			return ctx.Err()
		}

		return nil
	}
}
