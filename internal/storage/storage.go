package storagedb

import (
	"context"
	"database/sql"
	"fmt"
)

type Storage struct {
	DB  *sql.DB
	Ctx context.Context
}

func (m *Storage) Init(path string, ctx context.Context) (*sql.DB, error) {
	var err error
	m.Ctx = ctx
	m.DB, err = sql.Open("pgx", path)
	if err != nil {
		return nil, fmt.Errorf("fatal start a transaction %q", err)
	}

	err3 := m.Bootstrap(m.Ctx)
	if err3 != nil {
		return nil, fmt.Errorf("fatal start a transaction %q", err3)
	}
	return m.DB, nil

}

func (m *Storage) Bootstrap(ctx context.Context) error {
	// запускаем транзакцию
	tx, err := m.DB.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("fatal start a transaction %q", err)
	}
	// в случае неуспешного коммита все изменения транзакции будут отменены
	defer tx.Rollback()
	_, err = tx.ExecContext(m.Ctx, `
		CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			login varchar(100) UNIQUE NOT NULL, 
			password varchar(200) NOT NULL
		);
	`)
	if err != nil {
		return fmt.Errorf("fatal start a transaction %q", err)
	}
	_, err = tx.ExecContext(m.Ctx, `
		CREATE TABLE IF NOT EXISTS orders (
			id SERIAL PRIMARY KEY,
			login varchar(100) NOT NULL, 
			number bigint,
			status varchar(100),
			accrual DOUBLE PRECISION,
			uploadedat TIMESTAMP WITH TIME ZONE NOT NULL
		);
	`)
	if err != nil {
		return fmt.Errorf("fatal start a transaction %q", err)
	}
	_, err = tx.ExecContext(m.Ctx, `
		CREATE TABLE IF NOT EXISTS balances (
			id SERIAL PRIMARY KEY,
			login varchar(100) UNIQUE NOT NULL, 
			current DOUBLE PRECISION,
			withdrawn DOUBLE PRECISION
		);
	`)
	if err != nil {
		return fmt.Errorf("fatal start a transaction %q", err)
	}
	return tx.Commit()
}