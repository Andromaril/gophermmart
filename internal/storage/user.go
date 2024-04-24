package storagedb

import (
	"database/sql"
	"fmt"

	"github.com/andromaril/gophermmart/internal/errormart"
	log "github.com/sirupsen/logrus"
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

func (m *Storage) GetUser(login string) int {
	var value sql.NullInt64
	rows := m.DB.QueryRowContext(m.Ctx, "SELECT id FROM users WHERE login=$1", login)
	err := rows.Scan(&value)
	if err != nil {
		e := errormart.NewMartError(err)
		log.Error("error in select user id in user bd ", e.Error())
		return 0
	}
	if !value.Valid {
		log.Error("error in select in users bd: invalid user id")
		return 0
	}
	return int(value.Int64)
}

func (m *Storage) GetUserPassword(login string) string {
	var value sql.NullString
	rows := m.DB.QueryRowContext(m.Ctx, "SELECT password FROM users WHERE login=$1", login)
	err := rows.Scan(&value)
	if err != nil {
		e := errormart.NewMartError(err)
		log.Error("error in select password in user bd ", e.Error())
		return ""
	}
	if !value.Valid {
		log.Error("error in select in users bd: invalid password")
		return ""
	}
	return value.String
}
