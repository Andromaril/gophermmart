package main

import (
	"context"
	"database/sql"
	"net/http"
	"time"

	"github.com/go-chi/chi"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/andromaril/gophermmart/internal/accrual"
	"github.com/andromaril/gophermmart/internal/errormart"
	"github.com/andromaril/gophermmart/internal/flag"
	h "github.com/andromaril/gophermmart/internal/handler"
	"github.com/andromaril/gophermmart/internal/middleware"
	storagedb "github.com/andromaril/gophermmart/internal/storage"
	log "github.com/sirupsen/logrus"
)

func Update(newdb *storagedb.Storage) {
	for {
		log.Info("start send request to accrual")
		err := accrual.Accrual(newdb)

		if err != nil {
			e := errormart.NewMartError(err)
			log.Error("error witg accraul ", e.Error())
		}
		time.Sleep(time.Second * 5)
	}
}

func NewRouter(newdb storagedb.Storage) {
	r := chi.NewRouter()
	r.Post("/api/user/register", h.Register(newdb))
	r.Post("/api/user/login", h.Login(newdb))
	r.With(middleware.AuthMiddleware).Post("/api/user/orders", h.NewOrder(newdb))
	r.With(middleware.AuthMiddleware).Get("/api/user/orders", h.GetOrder(newdb))
	r.With(middleware.AuthMiddleware).Get("/api/user/balance", h.GetBalance(newdb))
	r.With(middleware.AuthMiddleware).Post("/api/user/balance/withdraw", h.NewWithdrawal(newdb))
	r.With(middleware.AuthMiddleware).Get("/api/user/withdrawals", h.GetWithdrawal(newdb))
	if err := http.ListenAndServe(flag.FlagRunAddr, r); err != nil {
		panic(err)

	}
	//return r
}

func main() {
	var err error
	flag.ParseFlags()
	var newdb storagedb.Storage
	var db *sql.DB
	db, err = newdb.Init(flag.Databaseflag, context.Background())
	if err != nil {
		panic(err)
	}
	defer db.Close()
	go Update(&newdb)
	// r := chi.NewRouter()
	// r.Post("/api/user/register", h.Register(newdb))
	// r.Post("/api/user/login", h.Login(newdb))
	// r.With(middleware.AuthMiddleware).Post("/api/user/orders", h.NewOrder(newdb))
	// r.With(middleware.AuthMiddleware).Get("/api/user/orders", h.GetOrder(newdb))
	// r.With(middleware.AuthMiddleware).Get("/api/user/balance", h.GetBalance(newdb))
	// r.With(middleware.AuthMiddleware).Post("/api/user/balance/withdraw", h.NewWithdrawal(newdb))
	// r.With(middleware.AuthMiddleware).Get("/api/user/withdrawals", h.GetWithdrawal(newdb))
	NewRouter(newdb)
	// if err := http.ListenAndServe(flag.FlagRunAddr, r); err != nil {
	// 	panic(err)

	// }
}
