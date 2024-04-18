package storagedb

import (
	"fmt"

	"github.com/andromaril/gophermmart/internal/model"
)

func (m *Storage) GetBalance(login string) (model.Balance, error) {
	result := model.Balance{}
	rows, err := m.DB.QueryContext(m.Ctx, "SELECT current, withdrawn FROM balances WHERE login=$1", login)
	//err := rows.Scan(&value)
	if err != nil {
		return model.Balance{}, fmt.Errorf("invalid login %q", err)
	}

	// обязательно закрываем перед возвратом функции
	defer rows.Close()
	rows.Next()
	err = rows.Scan(&result.Current, &result.Withdrawn)
	if err != nil {
		return model.Balance{}, fmt.Errorf("invalid login %q", err)
	}
	return result, nil
}
