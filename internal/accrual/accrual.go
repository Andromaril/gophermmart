package accrual

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/andromaril/gophermmart/internal/flag"
	"github.com/andromaril/gophermmart/internal/model"
	storagedb "github.com/andromaril/gophermmart/internal/storage"
	"github.com/go-resty/resty/v2"
)

func Accrual(storage storagedb.Storage) error {
	orders, err := storage.GetAccrualOrders()
	if err != nil {
		// f := fmt.Sprint("%q", err)
		// res.Write([]byte(f))
		// sugar.Errorw(
		// 	"error when get order")
		return err
	}

	for _, order := range orders {
		var updateorder model.Order
		client := resty.New()
		url := fmt.Sprintf("%s/api/orders/%s/", flag.BonusAddress, order.Number)
		//response, err2 := client.R().Get(client.BaseURL + "/api/orders/" + order.Number)
		response, err2 := client.R().Get(url)
		if err2 != nil {
			log.Println("1")
			log.Printf("%q", err2)
			return fmt.Errorf("error send request %q", err)
		}
		err3 := json.Unmarshal(response.Body(), &updateorder)
		if err3 != nil {
			log.Println(updateorder)
			log.Printf("%q", err3)
			return fmt.Errorf("error send request %q", err3)
		}

		if updateorder.Status != "REGISTERED" {
			err4 := storage.UpdateOrderAccrual(updateorder.Accrual, updateorder.Status, updateorder.Number)
			if err4 != nil {
				log.Println("4")
				log.Printf("%q", err4)
				return fmt.Errorf("error send request %q", err4)
			}
		}

		if updateorder.Accrual != nil {
			err := storage.UpdateBalanceAccrual(updateorder.Number, updateorder.Accrual)
			if err != nil {
				log.Println("5")
				log.Printf("%q", err)
				return fmt.Errorf("error send request %q", err)
			}
		}
	}
	return nil
}
