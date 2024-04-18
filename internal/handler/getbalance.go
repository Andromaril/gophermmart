package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	storagedb "github.com/andromaril/gophermmart/internal/storage"
)

func GetBalance(m storagedb.Storage) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		cookie, _ := req.Cookie("Login")
		res.Header().Set("Content-Type", "application/json")
		result, err := m.GetBalance(cookie.Value)
		if err != nil {
			f := fmt.Sprint("%q", err)
			res.Write([]byte(f))
			res.WriteHeader(http.StatusInternalServerError)
			return
		}
		body, err1 := json.Marshal(result)
		if err1 != nil {
			// f := fmt.Sprint("%q", err)
			// res.Write([]byte(f))
			res.WriteHeader(http.StatusInternalServerError)
			return
		}
		res.Write([]byte(body))
		res.WriteHeader(http.StatusOK)
	}
}
