package handler

import (
	"encoding/json"
	"net/http"

	"github.com/andromaril/gophermmart/internal/errormart"
	storagedb "github.com/andromaril/gophermmart/internal/storage"
	log "github.com/sirupsen/logrus"
)

func GetBalance(m storagedb.Storage) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		cookie, _ := req.Cookie("Login")
		res.Header().Set("Content-Type", "application/json")
		result, err := m.GetBalance(cookie.Value)
		log.Info("Current from user now: ", result.Current)
		if err != nil {
			e := errormart.NewMartError(err)
			log.Error(e.Error())
			res.WriteHeader(http.StatusInternalServerError)
			return
		}
		body, err1 := json.Marshal(result)
		if err1 != nil {
			e := errormart.NewMartError(err)
			log.Error(e.Error())
			res.WriteHeader(http.StatusInternalServerError)
			return
		}
		res.Write([]byte(body))
	}
}
