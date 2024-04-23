package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/andromaril/gophermmart/internal/model"
	storagedb "github.com/andromaril/gophermmart/internal/storage"
)

func GetOrder(m storagedb.Storage) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		r := make([]model.Order, 0)
		cookie, _ := req.Cookie("Login")
		res.Header().Set("Content-Type", "application/json")
		result, err := m.GetAllOrders(cookie.Value)
		if err != nil {
			// f := fmt.Sprint("%q", err)
			// res.Write([]byte(f))
			res.WriteHeader(http.StatusInternalServerError)
			return
		}
		for _, value := range result {
			r = append(r, model.Order{Number: value.Number, Status: value.Status, Accrual: value.Accrual, UploadedAt: value.UploadedAt})
			//log.Println(r)

		}
		if len(r) == 0 {
			res.WriteHeader(http.StatusNoContent)
		}
		enc := json.NewEncoder(res)
		if err := enc.Encode(r); err != nil {
			return
		}
		log.Println(r, cookie.Value)
		res.WriteHeader(http.StatusOK)
	}
}
