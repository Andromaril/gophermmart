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

func Register(m storagedb.Storage) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		var user model.User
		res.Header().Set("Content-Type", "application/json")
		dec := json.NewDecoder(req.Body)
		if err := dec.Decode(&user); err != nil {
			e := errormart.NewMartError(err)
			log.Error("error in decode request body from register ", e.Error())
			res.WriteHeader(http.StatusBadRequest)
			return
		}
		if user.Login == "" || user.Password == "" {
			log.Error("invalid registration data, empty login or password")
			res.WriteHeader(http.StatusBadRequest)
			return
		}
		hash := md5.Sum([]byte(user.Password))
		hashedPass := hex.EncodeToString(hash[:])
		value, err2 := m.GetUser(user.Login)
		if err2 != nil {
			e := errormart.NewMartError(err2)
			log.Error("error in select user id from users bd ", e.Error())
			res.WriteHeader(http.StatusInternalServerError)
			return
		}
		if value != "" {
			res.WriteHeader(http.StatusConflict)
			return
		}
		err3 := m.NewUser(user.Login, hashedPass)
		if err3 != nil {
			e := errormart.NewMartError(err3)
			log.Error("error in insert new user into users bd ", e.Error())
			res.WriteHeader(http.StatusInternalServerError)
			return
		}
		token, _ := verification.BuildJWTString()
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
