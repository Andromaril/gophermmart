package login

import (
	"encoding/json"
	"net/http"

	"github.com/andromaril/gophermmart/internal/model"
	storagedb "github.com/andromaril/gophermmart/internal/storage"
)

func Register(m storagedb.Storage) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		var user model.User
		res.Header().Set("Content-Type", "application/json")
		dec := json.NewDecoder(req.Body)
		if err := dec.Decode(&user); err != nil {
			res.WriteHeader(http.StatusBadRequest)
			return
		}
		if user.Login == "" || user.Password == "" {
			res.WriteHeader(http.StatusBadRequest)
			return
		}
		// value, _ := m.GetUser(user.Login)
		// if value != "" {
		// 	res.WriteHeader(http.StatusConflict)
		// 	return
		// }

		// err := m.NewUser(user.Login, user.Password)
		// if err != nil {
		// 	res.WriteHeader(http.StatusInternalServerError)
		// 	return
		// }
		resp := model.User{
			Login:    user.Login,
			Password: user.Password,
		}
		enc := json.NewEncoder(res)
		if err := enc.Encode(resp); err != nil {
			return
		}
		res.WriteHeader(http.StatusOK)
	}
}
