package storagedb

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/andromaril/gophermmart/internal/model"
)

var ErrNotBalance = errors.New("malo!!!")

func (m *Storage) GetWithdrawal(login string) ([]model.Withdrawn, error) {
	result := make([]model.Withdrawn, 0)
	rows, err := m.DB.QueryContext(m.Ctx, "SELECT orders.number, sum, processed_at FROM withdrawals INNER JOIN orders ON withdrawals.order_id = order_id WHERE withdrawals.login=$1", login)
	//err := rows.Scan(&value)
	if err != nil {
		log.Printf("%q", err)
		return result, fmt.Errorf("invalid login %q", err)
	}

	// обязательно закрываем перед возвратом функции
	defer rows.Close()
	for rows.Next() {
		var (
			order       string
			sum         float64
			processedat time.Time
		)
		err = rows.Scan(&order, &sum, &processedat)
		if err != nil {
			return result, fmt.Errorf("invalid login %q", err)
		}

		result = append(result, model.Withdrawn{Order: order, Sum: sum, ProcessedAt: processedat})
	}
	return result, nil
}

func (m *Storage) UpdateBalance(login string, withdrawal model.Withdrawn) error {

	balance, err := m.GetBalance(login)
	if err != nil {
		return fmt.Errorf("error insert %q", err)
	}
	number, _ := strconv.Atoi(withdrawal.Order)
	orderid, err1 := m.getOrderId(number)
	if err1 != nil {
		return fmt.Errorf("error insert %q", err)
	}
	user, _ := m.getUserId(login)
	if user != 0 {
		if balance.Current < withdrawal.Sum {
			return ErrNotBalance
		}
		balance2 := model.Balance{
			Current:   balance.Current - withdrawal.Sum,
			Withdrawn: balance.Withdrawn + withdrawal.Sum,
		}

		_, err2 := m.DB.ExecContext(m.Ctx, `
		UPDATE balances SET current=$1, withdrawn=$2 WHERE login=$3`, balance2.Current, balance2.Withdrawn, login)
		if err != nil {
			return fmt.Errorf("error insert %q", err2)
		}
	} else {

		balance, err := m.GetAccural(number)
		if err != nil {
			return fmt.Errorf("error insert %q", err)
		}
		_, err2 := m.DB.ExecContext(m.Ctx, `
		INSERT INTO balances (login, current, withdrawn) values ($1)`, login, balance, 0)
		if err != nil {
			return fmt.Errorf("error insert %q", err2)
		}
	}

	_, err3 := m.DB.ExecContext(m.Ctx, `
	INSERT INTO withdrawals (login, order_id, sum, processed_at)
	VALUES($1, $2, $3, $4)`, login, orderid, withdrawal.Sum, withdrawal.ProcessedAt)
	if err != nil {
		return fmt.Errorf("error insert %q", err3)
	}
	return nil
}

func (m *Storage) getUserId(login string) (int, error) {
	var value sql.NullInt64
	row := m.DB.QueryRowContext(m.Ctx, "SELECT id FROM balances WHERE login = $1", login)
	err := row.Scan(&value)
	if err != nil {
		return 0, fmt.Errorf("error select %q", err)
	}
	if !value.Valid {
		return 0, fmt.Errorf("invalid login %q", err)
	}

	return int(value.Int64), nil
}