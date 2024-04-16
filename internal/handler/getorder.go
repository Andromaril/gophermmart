package handler

import (
	"encoding/json"
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
			res.WriteHeader(http.StatusNoContent)
			return
		}
		for _, value := range result {
			r = append(r, model.Order{Number: value.Number, Status: value.Status, Accrual: value.Accrual, UploadedAt: value.UploadedAt})

		}
		enc := json.NewEncoder(res)
		if err := enc.Encode(r); err != nil {
			return
		}
		res.WriteHeader(http.StatusOK)
	}
}
