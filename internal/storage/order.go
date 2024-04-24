package storagedb

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/andromaril/gophermmart/internal/errormart"
	"github.com/andromaril/gophermmart/internal/model"
	log "github.com/sirupsen/logrus"
)

func (m *Storage) NewOrder(login string, order string) error {
	_, err := m.DB.ExecContext(m.Ctx, `
	INSERT INTO orders (login, number, status, uploadedat)
	VALUES($1, $2, $3, $4)`, login, order, "NEW", time.Now().Format(time.RFC3339))
	if err != nil {
		e := errormart.NewMartError(err)
		return fmt.Errorf("error insert %q", e.Error())
	}
	return nil
}

func (m *Storage) GetOrderUser(login string, order string) int {
	var value sql.NullInt64
	rows := m.DB.QueryRowContext(m.Ctx, "SELECT id FROM orders WHERE login=$1 AND number=$2", login, order)
	err := rows.Scan(&value)
	if err != nil {
		e := errormart.NewMartError(err)
		log.Error(e.Error())
		return 0
	}
	if !value.Valid {
		log.Error("invalid id")
		return 0
	}

	return int(value.Int64)
}

func (m *Storage) GetOrderAnotherUser(order string) string {
	var value sql.NullString

	rows := m.DB.QueryRowContext(m.Ctx, "SELECT login FROM orders WHERE number=$1", order)
	err := rows.Scan(&value)
	if err != nil {
		e := errormart.NewMartError(err)
		log.Error("error in select in orders db ", e.Error())
		return ""
	}
	if !value.Valid {
		log.Error("error in select in orders db: invalid login")
		return ""
	}
	return value.String
}

func (m *Storage) GetAllOrders(login string) ([]model.Order, error) {
	result := make([]model.Order, 0)
	rows, err := m.DB.QueryContext(m.Ctx, "SELECT number, status, accrual, uploadedat FROM orders WHERE login=$1", login)
	if err != nil {
		e := errormart.NewMartError(err)
		return result, fmt.Errorf("error select %q", e.Error())
	}
	defer rows.Close()
	for rows.Next() {
		var (
			number     string
			status     string
			accrual    *float64
			uploadedat time.Time
		)
		err = rows.Scan(&number, &status, &accrual, &uploadedat)
		if err != nil {
			e := errormart.NewMartError(err)
			return result, fmt.Errorf("invalid scan %q", e.Error())
		}
		result = append(result, model.Order{Number: number, Status: status, Accrual: accrual, UploadedAt: uploadedat})
	}
	err = rows.Err()
	if err != nil {
		e := errormart.NewMartError(err)
		return result, fmt.Errorf("error select %q", e.Error())
	}
	return result, nil
}

// func (m *Storage) getOrderId(number string) (int, error) {
// 	var value sql.NullInt64
// 	row := m.DB.QueryRowContext(m.Ctx, "SELECT id FROM orders WHERE number = $1", number)
// 	err := row.Scan(&value)
// 	if err != nil {
// 		return 0, fmt.Errorf("error select %q", err)
// 	}
// 	if !value.Valid {
// 		return 0, fmt.Errorf("invalid login %q", err)
// 	}

// 	return int(value.Int64), nil
// }

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

// func (m *Storage) GetAccural(login string) (float64, error) {
// 	var value sql.NullFloat64
// 	row := m.DB.QueryRowContext(m.Ctx, "SELECT accrual FROM orders WHERE login = $1", login)
// 	err := row.Scan(&value)
// 	if err != nil {
// 		return 0, fmt.Errorf("error select %q", err)
// 	}
// 	if !value.Valid {
// 		return 0, fmt.Errorf("invalid login %q", err)
// 	}

// 	return value.Float64, nil
// }

func (m *Storage) GetAccrualOrders() ([]model.Order, error) {
	result := make([]model.Order, 0)
	rows, err := m.DB.QueryContext(m.Ctx, "SELECT number, status, accrual, uploadedat FROM orders WHERE status = 'NEW' OR status = 'PROCESSING'")
	if err != nil {
		e := errormart.NewMartError(err)
		return result, fmt.Errorf("error select %q", e.Error())
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
			e := errormart.NewMartError(err)
			return result, fmt.Errorf("invalid scan %q", e.Error())
		}
		result = append(result, model.Order{Number: number, Status: status, Accrual: accrual, UploadedAt: uploadedat})
	}
	err = rows.Err()
	if err != nil {
		e := errormart.NewMartError(err)
		return result, fmt.Errorf("error select %q", e.Error())
	}
	return result, nil
}

func (m *Storage) UpdateOrderAccrual(accrual *float64, status string, number string) error {
	_, err := m.DB.ExecContext(m.Ctx, `
		UPDATE orders SET accrual=$1, status=$2 WHERE number=$3`, accrual, status, number)
	if err != nil {
		e := errormart.NewMartError(err)
		return fmt.Errorf("error insert %q", e.Error())
	}
	return nil
}

func (m *Storage) GetUserLogin(number string) (string, error) {
	var value sql.NullString
	row := m.DB.QueryRowContext(m.Ctx, "SELECT login FROM orders WHERE number = $1", number)
	err := row.Scan(&value)
	if err != nil {
		e := errormart.NewMartError(err)
		return "", fmt.Errorf("error select %q", e.Error())
	}
	if !value.Valid {
		e := errormart.NewMartError(err)
		return "", fmt.Errorf("invalid login %q", e.Error())
	}

	return value.String, nil
}
