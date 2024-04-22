package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/andromaril/gophermmart/internal/model"
	storagedb "github.com/andromaril/gophermmart/internal/storage"
	"github.com/theplant/luhn"
)

func GetWithdrawal(m storagedb.Storage) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		r := make([]model.Withdrawn, 0)
		cookie, _ := req.Cookie("Login")
		res.Header().Set("Content-Type", "application/json")
		result, err := m.GetWithdrawal(cookie.Value)
		if err != nil {
			// f := fmt.Sprint("%q", err)
			// res.Write([]byte(f))
			res.WriteHeader(http.StatusInternalServerError)
			return
		}
		for _, value := range result {
			r = append(r, model.Withdrawn{Order: value.Order, Sum: value.Sum, ProcessedAt: value.ProcessedAt})

		}
		if len(r) == 0 {
			res.WriteHeader(http.StatusNoContent)
		}
		enc := json.NewEncoder(res)
		if err := enc.Encode(r); err != nil {
			return
		}
		res.WriteHeader(http.StatusOK)
	}
}

func NewWithdrawal(m storagedb.Storage) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		var r model.Withdrawn
		dec := json.NewDecoder(req.Body)
		if err := dec.Decode(&r); err != nil {
			res.WriteHeader(http.StatusBadRequest)
			return
		}
		cookie, _ := req.Cookie("Login")
		number, _ := strconv.Atoi(r.Order)
		validnumer := luhn.Valid(number)
		// _, err1 := m.GetOrderUser(cookie.Value, int(number))
		// if err1 != nil {
		// 	//res.Write([]byte(cookie.Value))
		// 	log.Printf("%q", err1)
		// 	res.WriteHeader(http.StatusUnprocessableEntity)
		// 	return
		// }
		if validnumer {
			err := m.UpdateBalance(cookie.Value, r)
			if err != nil {
				// f := fmt.Sprint("%q", err)
				// res.Write([]byte(f))
				if err == storagedb.ErrNotBalance {
					res.WriteHeader(http.StatusPaymentRequired)
					return
				}
				log.Printf("%q", err)
				res.WriteHeader(http.StatusInternalServerError)
				return
			}
		} else {
			res.WriteHeader(http.StatusUnprocessableEntity)
		}
		res.WriteHeader(http.StatusOK)
	}
}
