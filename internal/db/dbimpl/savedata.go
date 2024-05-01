package dbimpl

import (
	"context"

	"github.com/gleb-korostelev/gophermart.git/internal/config"
	"github.com/gleb-korostelev/gophermart.git/internal/db"
	"github.com/gleb-korostelev/gophermart.git/internal/models"
	"github.com/jackc/pgx/v5"
)

func SaveUser(db db.DB, ctx context.Context, login, password string) error {
	sql := `
	INSERT INTO user_data (login, password, is_deleted)
	VALUES ($1, $2, FALSE)
	ON CONFLICT (login)
	DO UPDATE SET
		login = EXCLUDED.login,
		is_deleted = FALSE
	WHERE user_data.is_deleted = TRUE
	`

	sqlBalance := `
	INSERT INTO balances (login, current, withdrawn)
	VALUES ($1, 0, 0)`

	cmdTag, err := db.Exec(ctx, sql, login, password)
	if err != nil {
		return err
	}

	_, err = db.Exec(ctx, sqlBalance, login)
	if err != nil {
		return err
	}
	if cmdTag.RowsAffected() == 0 {
		return config.ErrLoginExists
	}
	return nil
}

func SaveOrders(db db.DB, ctx context.Context, login, orderID string) (string, bool, error) {
	var orderLogin string
	sqlScan := `SELECT login FROM orders WHERE order_id=$1`
	sqlExec := `INSERT INTO orders (order_id, login, status) VALUES ($1, $2, $3)`
	err := db.QueryRow(ctx, sqlScan, orderID).Scan(&orderLogin)
	if err != nil {
		if err == pgx.ErrNoRows {
			_, err = db.Exec(ctx, sqlExec, orderID, login, "NEW")
			if err != nil {
				return "", false, err
			}
			return "", false, nil
		}
		return "", false, err
	}
	return orderLogin, true, nil
}

func Withdraw(db db.DB, ctx context.Context, login string, req models.WithdrawRequest) error {
	tx, err := db.BeginW(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	var currentBalance float64
	err = tx.QueryRow(ctx, "SELECT current FROM balances WHERE login=$1", login).Scan(&currentBalance)
	if err != nil {
		return err
	}

	if currentBalance < req.Sum {
		return config.ErrNoFunds
	}

	_, err = tx.Exec(ctx, "UPDATE balances SET current = current - $1, withdrawn=withdrawn + $1 WHERE login = $2", req.Sum, login)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, "INSERT INTO withdrawals (login, order_id, sum, processed_at) VALUES ($1, $2, $3, NOW())", login, req.Order, req.Sum)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func UpdateOrderInfo(db db.DB, ctx context.Context, login string, orderInfo models.OrderResponse) error {
	sql := `UPDATE orders SET status = $1, accrual = COALESCE($2, accrual) WHERE login = $3`

	tx, err := db.BeginW(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, sql, orderInfo.Status, orderInfo.Accrual, login)
	if err != nil {
		return err
	}
	_, err = tx.Exec(ctx, "UPDATE balances SET current = current + $1 WHERE login = $2", orderInfo.Accrual, login)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}
