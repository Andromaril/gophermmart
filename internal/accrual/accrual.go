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
	orders, err := storage.GetAccrualOrders()
	if err != nil {
		e := errormart.NewMartError(err)
		log.Error(e.Error())
		return fmt.Errorf("error %q", e.Error())
	}

	for _, order := range orders {
		var updateorder model.UpdateOrder
		client := resty.New()
		url := fmt.Sprintf("%s/api/orders/%s", flag.BonusAddress, order.Number)
		response, err2 := client.R().Get(url)
		log.Info(response)
		if err2 != nil {
			e := errormart.NewMartError(err)
			log.Error(e.Error())
			return fmt.Errorf("error %q", e.Error())
		}
		err3 := json.Unmarshal(response.Body(), &updateorder)
		if err3 != nil {
			e := errormart.NewMartError(err)
			log.Error(e.Error())
			return fmt.Errorf("error %q", e.Error())
		}

		if updateorder.Status != "REGISTERED" {
			err4 := storage.UpdateOrderAccrual(updateorder.Accrual, updateorder.Status, updateorder.Number)
			if err4 != nil {
				e := errormart.NewMartError(err)
				log.Error(e.Error())
				return fmt.Errorf("error %q", e.Error())
			}
		}
		if updateorder.Accrual != nil {
			err5 := storage.UpdateBalanceAccrual(updateorder.Number, updateorder.Accrual)
			log.Info("update order wuth number", updateorder.Number)
			if err5 != nil {
				e := errormart.NewMartError(err5)
				log.Error(e.Error())
				return fmt.Errorf("error %q", e.Error())
			}
		}
	}
	return nil
}
