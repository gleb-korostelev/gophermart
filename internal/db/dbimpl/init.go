package dbimpl

import (
	"context"

	"github.com/gleb-korostelev/gophermart.git/internal/db"
)

func InitializeTables(db db.DB) error {
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS user_data (
	    id SERIAL PRIMARY KEY,
	    login VARCHAR(255) UNIQUE NOT NULL,
	    password VARCHAR(255) NOT NULL,
	    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
		is_deleted BOOLEAN DEFAULT FALSE
	);

	CREATE TABLE IF NOT EXISTS orders (
        id SERIAL PRIMARY KEY,
        login VARCHAR(255) NOT NULL,
        order_id VARCHAR(255) NOT NULL UNIQUE,
		status VARCHAR(255) NOT NULL,
		accrual FLOAT,
        uploaded_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
    );

	CREATE TABLE IF NOT EXISTS balances (
        id SERIAL PRIMARY KEY,
        login VARCHAR(255) NOT NULL,
        current_balance FLOAT NOT NULL,
		withdrawn FLOAT NOT NULL DEFAULT '0',
        uploaded_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
    );

	CREATE TABLE IF NOT EXISTS withdrawals  (
        id SERIAL PRIMARY KEY,
        login VARCHAR(255) NOT NULL,
        order_id VARCHAR(255) NOT NULL UNIQUE,
		sum FLOAT NOT NULL,
        processed_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
    );
	`
	_, err := db.Exec(context.Background(), createTableSQL)

	return err
}
