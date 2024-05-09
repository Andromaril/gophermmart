package storagedb

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/andromaril/gophermmart/internal/errormart"
	"github.com/andromaril/gophermmart/internal/model"
	log "github.com/sirupsen/logrus"
)

var ErrNotRow = errors.New("no balances, new user")

func (m *Storage) GetBalance(login string) (model.Balance, error) {
	result := model.Balance{}
	var current sql.NullFloat64
	var withdrawn sql.NullFloat64
	rows := m.DB.QueryRowContext(m.Ctx, "SELECT current, withdrawn FROM balances WHERE login=$1", login)
	err := rows.Scan(&current, &withdrawn)

	if errors.Is(err, sql.ErrNoRows) {
		return model.Balance{Current: 0, Withdrawn: 0}, ErrNotRow
	} else if err != nil {
		e := errormart.NewMartError(err)
		return model.Balance{}, fmt.Errorf("error select %q", e.Error())
	}

	if !current.Valid {
		log.Error("error in select in orders db: invalid login")
		return model.Balance{}, fmt.Errorf("error in select in balances db: invalid current")
	}
	if !withdrawn.Valid {
		log.Error("error in select in orders db: invalid login")
		return model.Balance{}, fmt.Errorf("error in select in balances db: invalid withdrawn")
	}

	result = model.Balance{Current: current.Float64, Withdrawn: withdrawn.Float64}
	return result, nil
}

func (m *Storage) UpdateBalanceAccrual(number string, accrual *float64) error {
	var err error
	login, err := m.GetUserLogin(number)
	if err != nil {
		e := errormart.NewMartError(err)
		return fmt.Errorf("error select %q", e.Error())
	}

	result, err := m.GetBalance(login)
	if err != ErrNotRow {
		e := errormart.NewMartError(err)
		return fmt.Errorf("error select balance %q", e.Error())
	}
	balancenew := model.Balance{
		Current:   result.Current + *accrual,
		Withdrawn: result.Withdrawn,
	}

	_, err = m.DB.ExecContext(m.Ctx, `
	 INSERT INTO balances (login, current, withdrawn) VALUES($1, $2, $3) ON CONFLICT (login) DO UPDATE SET current=$2, withdrawn=$3`, login, balancenew.Current, balancenew.Withdrawn)
	if err != nil {
		e := errormart.NewMartError(err)
		return fmt.Errorf("error insert %q", e.Error())
	}
	return nil
}
