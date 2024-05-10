package accrual

import (
	"encoding/json"
	"fmt"

	"github.com/andromaril/gophermmart/internal/errormart"
	"github.com/andromaril/gophermmart/internal/flag"
	"github.com/andromaril/gophermmart/internal/model"
	storagedb "github.com/andromaril/gophermmart/internal/storage"
	"github.com/go-resty/resty/v2"
	log "github.com/sirupsen/logrus"
)

func Accrual(storage *storagedb.Storage) error {
	var err error
	orders, err := storage.GetAccrualOrders()
	if err != nil {
		e := errormart.NewMartError(err)
		log.Error(e.Error())
		return fmt.Errorf("error %q", e.Error())
	}
	client := resty.New()
	for _, order := range orders {
		var updateorder model.UpdateOrder
		url := fmt.Sprintf("%s/api/orders/%s", flag.BonusAddress, order.Number)
		response, err := client.R().Get(url)
		log.Info(response)
		if err != nil {
			e := errormart.NewMartError(err)
			log.Error(e.Error())
			return fmt.Errorf("error %q", e.Error())
		}
		err = json.Unmarshal(response.Body(), &updateorder)
		if err != nil {
			e := errormart.NewMartError(err)
			log.Error(e.Error())
			return fmt.Errorf("error %q", e.Error())
		}

		if updateorder.Status != "REGISTERED" {
			err = storage.UpdateOrderAccrual(updateorder.Accrual, updateorder.Status, updateorder.Number)
			if err != nil {
				e := errormart.NewMartError(err)
				log.Error(e.Error())
				return fmt.Errorf("error %q", e.Error())
			}
		}
		if updateorder.Accrual != nil {
			err = storage.UpdateBalanceAccrual(updateorder.Number, updateorder.Accrual)
			log.Info("update order from accrual with number ", updateorder.Number)
			if err != nil {
				e := errormart.NewMartError(err)
				log.Error(e.Error())
				return fmt.Errorf("error %q", e.Error())
			}
		}
	}
	return nil
}
