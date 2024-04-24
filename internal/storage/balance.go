package storagedb

import (
	"fmt"

	"github.com/andromaril/gophermmart/internal/errormart"
	"github.com/andromaril/gophermmart/internal/model"
)

func (m *Storage) GetBalance(login string) (model.Balance, error) {
	result := model.Balance{}
	rows, err := m.DB.QueryContext(m.Ctx, "SELECT current, withdrawn FROM balances WHERE login=$1", login)
	if err != nil {
		e := errormart.NewMartError(err)
		return model.Balance{}, fmt.Errorf("error select %q", e.Error())
	}

	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&result.Current, &result.Withdrawn)
		if err != nil {
			e := errormart.NewMartError(err)
			return model.Balance{}, fmt.Errorf("error select %q", e.Error())
		}
	}
	err = rows.Err()
	if err != nil {
		e := errormart.NewMartError(err)
		return model.Balance{}, fmt.Errorf("error select %q", e.Error())
	}
	return result, nil
}

func (m *Storage) UpdateBalanceAccrual(number string, accrual *float64) error {
	login, err := m.GetUserLogin(number)
	if err != nil {
		e := errormart.NewMartError(err)
		return fmt.Errorf("error select %q", e.Error())
	}

	result, err := m.GetBalance(login)
	if err != nil {
		e := errormart.NewMartError(err)
		return fmt.Errorf("error select %q", e.Error())
	}
	balancenew := model.Balance{
		Current:   result.Current + *accrual,
		Withdrawn: result.Withdrawn,
	}

	_, err2 := m.DB.ExecContext(m.Ctx, `
	 INSERT INTO balances (login, current, withdrawn) VALUES($1, $2, $3) ON CONFLICT (login) DO UPDATE SET current=$2, withdrawn=$3`, login, balancenew.Current, balancenew.Withdrawn)
	if err2 != nil {
		e := errormart.NewMartError(err)
		return fmt.Errorf("error insert %q", e.Error())
	}
	return nil
}
