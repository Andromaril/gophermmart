package storagedb

import (
	"fmt"
	"log"

	"github.com/andromaril/gophermmart/internal/model"
)

func (m *Storage) GetBalance(login string) (model.Balance, error) {
	result := model.Balance{}
	rows, err := m.DB.QueryContext(m.Ctx, "SELECT current, withdrawn FROM balances WHERE login=$1", login)
	//err := rows.Scan(&value)
	if err != nil {
		log.Println("This is a log message!")
		return model.Balance{}, fmt.Errorf("invalid login %q", err)
	}

	// обязательно закрываем перед возвратом функции
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&result.Current, &result.Withdrawn)
		if err != nil {
			return model.Balance{}, fmt.Errorf("invalid login %q", err)
		}
	}
	return result, nil
}

func (m *Storage) UpdateBalanceAccrual(number string, accrual *float64) error {
	login, err := m.GetUserLogin(number)
	if err != nil {
		log.Println("This is not login!")
		return fmt.Errorf("invalid login %q", err)
	}

	result, err := m.GetBalance(login)
	if err != nil {
		log.Println("This is not balance!")
		return fmt.Errorf("invalid balance %q", err)
	}
	balancenew := model.Balance{
		Current:   result.Current + *accrual,
		Withdrawn: result.Withdrawn,
	}
	_, err2 := m.DB.ExecContext(m.Ctx, `
	UPDATE balances SET current=$1 WHERE login=$2`, balancenew.Current, login)
	if err2 != nil {
		return fmt.Errorf("error insert3 %q", err2)
	}
	return nil
}
