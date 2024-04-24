package handler

import (
	"encoding/json"
	"net/http"

	"github.com/andromaril/gophermmart/internal/errormart"
	"github.com/andromaril/gophermmart/internal/model"
	storagedb "github.com/andromaril/gophermmart/internal/storage"
	"github.com/andromaril/gophermmart/internal/utils"
	log "github.com/sirupsen/logrus"
)

func GetWithdrawal(m storagedb.Storage) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		r := make([]model.Withdrawn, 0)
		cookie, _ := req.Cookie("Login")
		res.Header().Set("Content-Type", "application/json")
		result, err := m.GetWithdrawal(cookie.Value)
		if err != nil {
			e := errormart.NewMartError(err)
			log.Error("error in select number, sum, processed_at from withdrawals bd ", e.Error())
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
			e := errormart.NewMartError(err)
			log.Error("error in encode model.Withdrawn ", e.Error())
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
			e := errormart.NewMartError(err)
			log.Error("error in decode request body from withdrawal ", e.Error())
			res.WriteHeader(http.StatusBadRequest)
			return
		}
		cookie, _ := req.Cookie("Login")
		validnumer, err2 := utils.ValidLuhn(r.Order)
		if err2 != nil {
			e := errormart.NewMartError(err2)
			log.Error("error in valid luhn order number ", e.Error())
			res.WriteHeader(http.StatusInternalServerError)
			return
		}
		if validnumer {
			err := m.UpdateBalance(cookie.Value, r)
			if err != nil {
				e := errormart.NewMartError(err2)
				log.Error("error update withdrawals and balances bd ", e.Error())
				if err == storagedb.ErrNotBalance {
					res.WriteHeader(http.StatusPaymentRequired)
					return
				}
				res.WriteHeader(http.StatusInternalServerError)
				return
			}
		} else {
			res.WriteHeader(http.StatusUnprocessableEntity)
		}
	}
}
