package handler

import (
	"io"
	"net/http"
	"strconv"

	"github.com/andromaril/gophermmart/internal/errormart"
	storagedb "github.com/andromaril/gophermmart/internal/storage"
	log "github.com/sirupsen/logrus"
	"github.com/theplant/luhn"
)

func NewOrder(m storagedb.Storage) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {

		res.Header().Set("Content-Type", "text/plain")
		cookie, _ := req.Cookie("Login")
		requestData, err1 := io.ReadAll(req.Body)
		if err1 != nil {
			//res.Write([]byte(requestData))
			res.WriteHeader(http.StatusInternalServerError)
			return
		}
		number := string(requestData)
		number2, _ := strconv.Atoi(number)
		validnumer := luhn.Valid(number2)
		if validnumer {
			orderexist, err3 := m.GetOrderUser(cookie.Value, number)
			if orderexist != 0 {
				e := errormart.NewMartError(err3)
				log.Error(e.Error())
				res.WriteHeader(http.StatusOK)
				return
			}
			orderexist2, _ := m.GetOrderAnotherUser(number)
			if orderexist2 != "" && orderexist2 != cookie.Value {
				res.WriteHeader(http.StatusConflict)
				return
			}
			err := m.NewOrder(cookie.Value, number)
			if err != nil {
				// f := fmt.Sprint("%q", err)
				// res.Write([]byte(f))
				res.WriteHeader(http.StatusInternalServerError)
				return
			}
		} else {
			res.WriteHeader(http.StatusUnprocessableEntity)
		}
		res.WriteHeader(http.StatusAccepted)
	}
}
