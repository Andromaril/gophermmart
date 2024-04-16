package storagedb

import (
	"database/sql"
	"fmt"
)

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
