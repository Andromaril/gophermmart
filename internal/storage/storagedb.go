package storagedb

import (
	"context"
	"database/sql"
	"fmt"
	"time"
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
			login varchar(100) UNIQUE NOT NULL, 
			number bigint,
			status varchar(100),
			accrual DOUBLE PRECISION,
			uploadedat TIMESTAMP WITH TIME ZONE NOT NULL
		);
	`)
	if err != nil {
		return fmt.Errorf("fatal start a transaction %q", err)
	}
	return tx.Commit()
}

func (m *Storage) NewUser(login string, password string) error {
	_, err := m.DB.ExecContext(m.Ctx, `
	INSERT INTO users (login, password)
	VALUES($1, $2)`, login, password)
	if err != nil {
		return fmt.Errorf("error insert %q", err)
	}
	return nil
}

func (m *Storage) GetUser(login string) (string, error) {
	var value sql.NullString
	rows := m.DB.QueryRowContext(m.Ctx, "SELECT id FROM users WHERE login=$1", login)
	err := rows.Scan(&value)
	if err != nil {
		return "", fmt.Errorf("error select %q", err)
	}
	if !value.Valid {
		return "", fmt.Errorf("invalid login %q", err)
	}
	return value.String, nil
}

func (m *Storage) GetUserPassword(login string) (string, error) {
	var value sql.NullString
	rows := m.DB.QueryRowContext(m.Ctx, "SELECT password FROM users WHERE login=$1", login)
	err := rows.Scan(&value)
	if err != nil {
		return "", fmt.Errorf("error select %q", err)
	}
	if !value.Valid {
		return "", fmt.Errorf("invalid login %q", err)
	}
	return value.String, nil
}

func (m *Storage) NewOrder(login string, order int) error {
	_, err := m.DB.ExecContext(m.Ctx, `
	INSERT INTO orders (login, number, status, uploadedat)
	VALUES($1, $2, $3, $4)`, login, order, "NEW", time.Now().Format(time.RFC3339))
	if err != nil {
		return fmt.Errorf("error insert %q", err)
	}
	return nil
}

func (m *Storage) GetOrderUser(login string, order int) (int, error) {
	var value sql.NullInt64
	rows := m.DB.QueryRowContext(m.Ctx, "SELECT id FROM orders WHERE login=$1 AND number=$2", login, order)
	err := rows.Scan(&value)
	if err != nil {
		return 0, fmt.Errorf("error select %q", err)
	}
	if !value.Valid {
		return 0, fmt.Errorf("invalid login %q", err)
	}

	return int(value.Int64), nil
}

func (m *Storage) GetOrderAnotherUser(order int) (string, error) {
	var value sql.NullString

	rows2 := m.DB.QueryRowContext(m.Ctx, "SELECT login FROM orders WHERE number=$1", order)
	err2 := rows2.Scan(&value)
	if err2 != nil {
		return "", fmt.Errorf("error select %q", err2)
	}
	if !value.Valid {
		return "", fmt.Errorf("invalid login %q", err2)
	}
	return value.String, nil
}