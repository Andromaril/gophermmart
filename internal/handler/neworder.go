package handler

import (
	"io"
	"net/http"
	"strconv"

	storagedb "github.com/andromaril/gophermmart/internal/storage"
	"github.com/theplant/luhn"
)

func NewOrder(m storagedb.Storage) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {

		res.Header().Set("Content-Type", "text/plain")
		cookie, _ := req.Cookie("Login")
		requestData, err1 := io.ReadAll(req.Body)
		if err1 != nil {
			res.WriteHeader(http.StatusInternalServerError)
			return
		}
		number := string(requestData)
		number2, _ := strconv.Atoi(number)
		validnumer := luhn.Valid(number2)
		if validnumer {
			orderexist, _ := m.GetOrderUser(cookie.Value, number2)
			if orderexist != 0 {
				res.Write([]byte(cookie.Value))
				res.WriteHeader(http.StatusOK)
				return
			}
			orderexist2, _ := m.GetOrderAnotherUser(number2)
			if orderexist2 != "" && orderexist2 != cookie.Value {
				res.WriteHeader(http.StatusConflict)
				return
			}
			err := m.NewOrder(cookie.Value, number2)
			if err != nil {
				res.WriteHeader(http.StatusInternalServerError)
				return
			}
		} else {
			res.WriteHeader(http.StatusUnprocessableEntity)
		}

	}
}
