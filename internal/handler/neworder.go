package handler

import (
	"io"
	"net/http"

	"github.com/andromaril/gophermmart/internal/errormart"
	storagedb "github.com/andromaril/gophermmart/internal/storage"
	"github.com/andromaril/gophermmart/internal/utils"
	log "github.com/sirupsen/logrus"
)

func NewOrder(m storagedb.Storage) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {

		res.Header().Set("Content-Type", "text/plain")
		cookie, _ := req.Cookie("Login")
		requestData, err1 := io.ReadAll(req.Body)
		if err1 != nil {
			e := errormart.NewMartError(err1)
			log.Error("error in read request data from order ", e.Error())
			res.WriteHeader(http.StatusInternalServerError)
			return
		}
		number := string(requestData)
		validnumer, err2 := utils.ValidLuhn(number)
		if err2 != nil {
			e := errormart.NewMartError(err2)
			log.Error("error in valid luhn order number ", e.Error())
			res.WriteHeader(http.StatusInternalServerError)
			return
		}
		if validnumer {
			orderexist := m.GetOrderUser(cookie.Value, number)
			if orderexist != 0 {
				res.WriteHeader(http.StatusOK)
				return
			}
			orderexist2 := m.GetOrderAnotherUser(number)
			if orderexist2 != "" && orderexist2 != cookie.Value {
				res.WriteHeader(http.StatusConflict)
				return
			}
			err := m.NewOrder(cookie.Value, number)
			if err != nil {
				e := errormart.NewMartError(err)
				log.Error("error in insert new order into orders bd ", e.Error())
				res.WriteHeader(http.StatusInternalServerError)
				return
			}
		} else {
			res.WriteHeader(http.StatusUnprocessableEntity)
		}
		res.WriteHeader(http.StatusAccepted)
	}
}
