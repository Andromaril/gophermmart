package handler

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"net/http"

	"github.com/andromaril/gophermmart/internal/model"
	storagedb "github.com/andromaril/gophermmart/internal/storage"
	"github.com/andromaril/gophermmart/internal/verification"
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
		hash := md5.Sum([]byte(user.Password))
		hashedPass := hex.EncodeToString(hash[:])
		value, _ := m.GetUser(user.Login)
		if value != "" {
			res.WriteHeader(http.StatusConflict)
			return
		}
		err := m.NewUser(user.Login, hashedPass)
		if err != nil {
			res.WriteHeader(http.StatusInternalServerError)
			return
		}
		token, _ := verification.BuildJWTString()
		cookie := &http.Cookie{
			Name:   "Token",
			Value:  token,
			MaxAge: 300,
		}
		http.SetCookie(res, cookie)
		res.WriteHeader(http.StatusOK)
	}
}
