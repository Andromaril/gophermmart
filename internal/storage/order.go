package storagedb

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/andromaril/gophermmart/internal/model"
)

func (m *Storage) NewOrder(login string, order string) error {
	_, err := m.DB.ExecContext(m.Ctx, `
	INSERT INTO orders (login, number, status, accrual, uploadedat)
	VALUES($1, $2, $3, $4)`, login, order, "NEW", 0, time.Now().Format(time.RFC3339))
	if err != nil {
		return fmt.Errorf("error insert %q", err)
	}
	return nil
}

func (m *Storage) GetOrderUser(login string, order string) (int, error) {
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

func (m *Storage) GetOrderAnotherUser(order string) (string, error) {
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

func (m *Storage) GetAllOrders(login string) ([]model.Order, error) {
	result := make([]model.Order, 0)
	rows, err := m.DB.QueryContext(m.Ctx, "SELECT number, status, accrual, uploadedat FROM orders WHERE login=$1", login)
	if err != nil {
		return result, fmt.Errorf("invalid login %q", err)
	}

	// обязательно закрываем перед возвратом функции
	defer rows.Close()

	// пробегаем по всем записям
	for rows.Next() {
		var (
			number     string
			status     string
			accrual    *float64
			uploadedat time.Time
		)
		err = rows.Scan(&number, &status, &accrual, &uploadedat)
		if err != nil {
			return result, fmt.Errorf("invalid login %q", err)
		}
		result = append(result, model.Order{Number: number, Status: status, Accrual: accrual, UploadedAt: uploadedat})
	}
	err = rows.Err()
	if err != nil {
		return result, fmt.Errorf("invalid login %q", err)
	}
	return result, nil
}

func (m *Storage) getOrderId(number string) (int, error) {
	var value sql.NullInt64
	row := m.DB.QueryRowContext(m.Ctx, "SELECT id FROM orders WHERE number = $1", number)
	err := row.Scan(&value)
	if err != nil {
		return 0, fmt.Errorf("error select %q", err)
	}
	if !value.Valid {
		return 0, fmt.Errorf("invalid login %q", err)
	}

	return int(value.Int64), nil
}

// func (m *Storage) GetAccural(number string) (float64, error) {
// 	var value sql.NullFloat64
// 	row := m.DB.QueryRowContext(m.Ctx, "SELECT accrual FROM orders WHERE number = $1", number)
// 	err := row.Scan(&value)
// 	if err != nil {
// 		return 0, fmt.Errorf("error select %q", err)
// 	}
// 	if !value.Valid {
// 		return 0, fmt.Errorf("invalid login %q", err)
// 	}

// 	return value.Float64, nil
// }

func (m *Storage) GetAccural(login string) (float64, error) {
	var value sql.NullFloat64
	row := m.DB.QueryRowContext(m.Ctx, "SELECT accrual FROM orders WHERE login = $1", login)
	err := row.Scan(&value)
	if err != nil {
		return 0, fmt.Errorf("error select %q", err)
	}
	if !value.Valid {
		return 0, fmt.Errorf("invalid login %q", err)
	}

	return value.Float64, nil
}
