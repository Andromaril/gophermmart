package storagedb

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/andromaril/gophermmart/internal/errormart"
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
		e := errormart.NewMartError(err)
		return nil, fmt.Errorf("fatal start a transaction %q", e.Error())
	}

	err = m.Bootstrap(m.Ctx)
	if err != nil {
		e := errormart.NewMartError(err)
		return nil, fmt.Errorf("fatal start a transaction %q", e.Error())
	}
	return m.DB, nil

}

func (m *Storage) Bootstrap(ctx context.Context) error {
	tx, err := m.DB.BeginTx(ctx, nil)
	if err != nil {
		e := errormart.NewMartError(err)
		return fmt.Errorf("fatal start a transaction %q", e.Error())
	}
	defer tx.Rollback()
	_, err = tx.ExecContext(m.Ctx, `
		CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			login varchar(100) UNIQUE NOT NULL, 
			password varchar(200) NOT NULL
		);
	`)
	if err != nil {
		e := errormart.NewMartError(err)
		return fmt.Errorf("fatal start a transaction %q", e.Error())
	}
	_, err = tx.ExecContext(m.Ctx, `
		CREATE TABLE IF NOT EXISTS orders (
			id SERIAL PRIMARY KEY,
			login varchar(100) NOT NULL, 
			number varchar(100) NOT NULL,
			status varchar(100),
			accrual DOUBLE PRECISION,
			uploadedat TIMESTAMP WITH TIME ZONE NOT NULL
		);
	`)
	if err != nil {
		e := errormart.NewMartError(err)
		return fmt.Errorf("fatal start a transaction %q", e.Error())
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
		e := errormart.NewMartError(err)
		return fmt.Errorf("fatal start a transaction %q", e.Error())
	}
	_, err = tx.ExecContext(m.Ctx, `
		CREATE TABLE IF NOT EXISTS withdrawals (
			id SERIAL PRIMARY KEY,
			number varchar(100) NOT NULL,
			login varchar(100) NOT NULL, 
			sum DOUBLE PRECISION,
			processed_at TIMESTAMP WITH TIME ZONE NOT NULL
		);
	`)
	if err != nil {
		e := errormart.NewMartError(err)
		return fmt.Errorf("fatal start a transaction %q", e.Error())
	}
	return tx.Commit()
}
