package storagedb

import (
	"errors"
	"fmt"
	"time"

	"github.com/andromaril/gophermmart/internal/errormart"
	"github.com/andromaril/gophermmart/internal/model"
)

var ErrNotBalance = errors.New("insufficient number of points to be deducted")

func (m *Storage) GetWithdrawal(login string) ([]model.Withdrawn, error) {
	result := make([]model.Withdrawn, 0)
	rows, err := m.DB.QueryContext(m.Ctx, "SELECT number, sum, processed_at FROM withdrawals WHERE login=$1", login)
	if err != nil {
		e := errormart.NewMartError(err)
		return result, fmt.Errorf("error select %q", e.Error())
	}

	defer rows.Close()
	for rows.Next() {
		var (
			order       string
			sum         float64
			processedat time.Time
		)
		err = rows.Scan(&order, &sum, &processedat)
		if err != nil {
			e := errormart.NewMartError(err)
			return result, fmt.Errorf("error scan %q", e.Error())
		}

		result = append(result, model.Withdrawn{Order: order, Sum: sum, ProcessedAt: processedat})
	}
	err = rows.Err()
	if err != nil {
		e := errormart.NewMartError(err)
		return result, fmt.Errorf("error select %q", e.Error())
	}
	return result, nil
}

func (m *Storage) UpdateBalance(login string, withdrawal model.Withdrawn) error {

	balance, err := m.GetBalance(login)
	if err != nil {
		e := errormart.NewMartError(err)
		return fmt.Errorf("error in select balance %q", e.Error())
	}
	if balance.Current < withdrawal.Sum {
		return ErrNotBalance
	}
	balance2 := model.Balance{
		Current:   balance.Current - withdrawal.Sum,
		Withdrawn: balance.Withdrawn + withdrawal.Sum,
	}

	_, err2 := m.DB.ExecContext(m.Ctx, `
		UPDATE balances SET current=$1, withdrawn=$2 WHERE login=$3`, balance2.Current, balance2.Withdrawn, login)
	if err2 != nil {
		e := errormart.NewMartError(err2)
		return fmt.Errorf("error insert %q", e.Error())
	}

	_, err3 := m.DB.ExecContext(m.Ctx, `
	INSERT INTO withdrawals (login, number, sum, processed_at)
	VALUES ($1, $2, $3, $4)`, login, withdrawal.Order, withdrawal.Sum, withdrawal.ProcessedAt)
	if err3 != nil {
		e := errormart.NewMartError(err3)
		return fmt.Errorf("error insert %q", e.Error())
	}
	return nil
}
