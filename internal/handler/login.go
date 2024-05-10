package handler

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"net/http"

	"github.com/andromaril/gophermmart/internal/errormart"
	"github.com/andromaril/gophermmart/internal/model"
	storagedb "github.com/andromaril/gophermmart/internal/storage"
	"github.com/andromaril/gophermmart/internal/verification"
	log "github.com/sirupsen/logrus"
)

func Login(m storagedb.Storage) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		var user model.User
		res.Header().Set("Content-Type", "application/json")
		dec := json.NewDecoder(req.Body)
		if err := dec.Decode(&user); err != nil {
			e := errormart.NewMartError(err)
			log.Error("error in decode request body from login ", e.Error())
			res.WriteHeader(http.StatusBadRequest)
			return
		}
		if user.Login == "" || user.Password == "" {
			log.Error("invalid login data, empty login or password")
			res.WriteHeader(http.StatusBadRequest)
			return
		}
		hash := md5.Sum([]byte(user.Password))
		hashedPass := hex.EncodeToString(hash[:])
		value := m.GetUserPassword(user.Login)
		if value == "" {
			res.WriteHeader(http.StatusInternalServerError)
			return
		}
		if value != hashedPass {
			res.WriteHeader(http.StatusUnauthorized)
			return
		}
		token, err := verification.BuildJWTString()
		if err != nil {
			e := errormart.NewMartError(err)
			log.Error("error in create token", e.Error())
			res.WriteHeader(http.StatusInternalServerError)
			return
		}
		cookie := &http.Cookie{
			Name:   "Token",
			Value:  token,
			Path:   "/",
			MaxAge: 300,
		}
		cookie2 := &http.Cookie{
			Name:   "Login",
			Value:  user.Login,
			Path:   "/",
			MaxAge: 300,
		}
		res.Header().Add("Authorization", user.Login)
		http.SetCookie(res, cookie)
		http.SetCookie(res, cookie2)
	}
}
