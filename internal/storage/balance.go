package storagedb

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/andromaril/gophermmart/internal/errormart"
	"github.com/andromaril/gophermmart/internal/model"
	"github.com/jackc/pgx/v4"
)

var ErrNotRow = errors.New("no balances, new user")

func (m *Storage) GetBalance(login string) (model.Balance, error) {
	result := model.Balance{}
	var current sql.NullFloat64
	var withdrawn sql.NullFloat64
	rows := m.DB.QueryRowContext(m.Ctx, "SELECT current, withdrawn FROM balances WHERE login=$1", login)
	// if err != nil {
	// 	e := errormart.NewMartError(err)
	// 	return model.Balance{}, fmt.Errorf("error select %q", e.Error())
	// }

	//defer rows.Close()
	//rows.Next()
	err := rows.Scan(current, withdrawn)
	// if err != nil {
	// 	e := errormart.NewMartError(err)
	// 	return model.Balance{}, fmt.Errorf("error select %q", e.Error())
	// }
	if errors.Is(err, pgx.ErrNoRows) {
		return model.Balance{0, 0}, ErrNotRow
	}
	result = model.Balance{current.Float64, withdrawn.Float64}
	//}
	// err = rows.Err()
	// if err != nil {
	// 	e := errormart.NewMartError(err)
	// 	return model.Balance{}, fmt.Errorf("error select %q", e.Error())
	// err := rows.Scan(&current, &withdrawn)
	// if err != nil {
	// 	e := errormart.NewMartError(err)
	// 	log.Error("error in select in orders db ", e.Error())
	// 	return 0, 0
	// }
	// if !current.Valid {
	// 	log.Error("error in select in orders db: invalid login")
	// 	return 0, 0
	// }
	// if !withdrawn.Valid {
	// 	log.Error("error in select in orders db: invalid login")
	// 	return 0,0
	// }
	return result, nil
	// }
	//return result, nil
}

func (m *Storage) UpdateBalanceAccrual(number string, accrual *float64) error {
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

	_, err2 := m.DB.ExecContext(m.Ctx, `
	 INSERT INTO balances (login, current, withdrawn) VALUES($1, $2, $3) ON CONFLICT (login) DO UPDATE SET current=$2, withdrawn=$3`, login, balancenew.Current, balancenew.Withdrawn)
	if err2 != nil {
		e := errormart.NewMartError(err)
		return fmt.Errorf("error insert %q", e.Error())
	}
	return nil
}
