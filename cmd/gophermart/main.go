package main

import (
	"context"
	"database/sql"
	"net/http"

	"github.com/go-chi/chi"
	"go.uber.org/zap"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/andromaril/gophermmart/internal/flag"
	h "github.com/andromaril/gophermmart/internal/handler"
	"github.com/andromaril/gophermmart/internal/middleware"
	storagedb "github.com/andromaril/gophermmart/internal/storage"
)

var sugar zap.SugaredLogger

func main() {
	var err error
	//var storage storagedb.Storage
	flag.ParseFlags()
	logger, err2 := zap.NewDevelopment()
	if err2 != nil {
		panic(err2)
	}
	defer logger.Sync()
	sugar = *logger.Sugar()
	sugar.Infow("Starting server")
	var newdb storagedb.Storage
	var db *sql.DB
	//if flag.Databaseflag != "" {
	db, err = newdb.Init(flag.Databaseflag, context.Background())
	if err != nil {
		panic(err)
	}
	defer db.Close()
	//}
	//defer db.Close()
	//}
	// var i int64
	// for i = 0; ; i++ {
	//time.Sleep(time.Second*5)
	// err1 := accrual.Accrual(&newdb, sugar)

	// if err1 != nil {
	// 	sugar.Infow("Not starting")
	// }
	// }
	r := chi.NewRouter()
	//r.Use(middleware.AuthMiddlewareContext)
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
	//var i int64
	//for i = 0; ; i++ {
	//time.Sleep(time.Second)

	//}
}
