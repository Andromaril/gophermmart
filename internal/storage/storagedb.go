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
			password varchar(100) UNIQUE NOT NULL
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
