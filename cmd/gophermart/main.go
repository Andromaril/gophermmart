package main

import (
	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/andromaril/gophermmart/internal/flag"
	"github.com/andromaril/gophermmart/internal/start"
)

// func Update(newdb *storagedb.Storage) {
// 	for {
// 		log.Info("start send request to accrual")
// 		err := accrual.Accrual(newdb)

// 		if err != nil {
// 			e := errormart.NewMartError(err)
// 			log.Error("error witg accraul ", e.Error())
// 		}
// 		time.Sleep(time.Second * 5)
// 	}
// }

// func NewRouter(newdb storagedb.Storage) {
// 	r := chi.NewRouter()
// 	r.Post("/api/user/register", h.Register(newdb))
// 	r.Post("/api/user/login", h.Login(newdb))
// 	r.With(middleware.AuthMiddleware).Post("/api/user/orders", h.NewOrder(newdb))
// 	r.With(middleware.AuthMiddleware).Get("/api/user/orders", h.GetOrder(newdb))
// 	r.With(middleware.AuthMiddleware).Get("/api/user/balance", h.GetBalance(newdb))
// 	r.With(middleware.AuthMiddleware).Post("/api/user/balance/withdraw", h.NewWithdrawal(newdb))
// 	r.With(middleware.AuthMiddleware).Get("/api/user/withdrawals", h.GetWithdrawal(newdb))
// 	if err := http.ListenAndServe(flag.FlagRunAddr, r); err != nil {
// 		panic(err)

// 	}
// 	//return r
// }
// func Start() (*sql.DB, storagedb.Storage) {
// 	var err error
// 	var newdb storagedb.Storage
// 	var db *sql.DB
// 	db, err = newdb.Init(flag.Databaseflag, context.Background())
// 	if err != nil {
// 		panic(err)
// 	}
// 	return db, newdb
// }

func main() {
	//var err error
	flag.ParseFlags()
	// var newdb storagedb.Storage
	// var db *sql.DB
	// db, err = newdb.Init(flag.Databaseflag, context.Background())
	// if err != nil {
	// 	panic(err)
	// }
	db, newdb := start.Start()
	defer db.Close()
	go start.Update(&newdb)
	// r := chi.NewRouter()
	// r.Post("/api/user/register", h.Register(newdb))
	// r.Post("/api/user/login", h.Login(newdb))
	// r.With(middleware.AuthMiddleware).Post("/api/user/orders", h.NewOrder(newdb))
	// r.With(middleware.AuthMiddleware).Get("/api/user/orders", h.GetOrder(newdb))
	// r.With(middleware.AuthMiddleware).Get("/api/user/balance", h.GetBalance(newdb))
	// r.With(middleware.AuthMiddleware).Post("/api/user/balance/withdraw", h.NewWithdrawal(newdb))
	// r.With(middleware.AuthMiddleware).Get("/api/user/withdrawals", h.GetWithdrawal(newdb))
	start.NewRouter(newdb)
	// if err := http.ListenAndServe(flag.FlagRunAddr, r); err != nil {
	// 	panic(err)

	// }
}
