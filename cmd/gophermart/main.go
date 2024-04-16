package main

import (
	"context"
	"database/sql"
	"net/http"

	"github.com/go-chi/chi"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/andromaril/gophermmart/internal/flag"
	h "github.com/andromaril/gophermmart/internal/handler"
	"github.com/andromaril/gophermmart/internal/middleware"
	storagedb "github.com/andromaril/gophermmart/internal/storage"
)

func main() {
	var err error
	//var storage storagedb.Storage
	flag.ParseFlags()
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
	r := chi.NewRouter()
	//r.Use(middleware.AuthMiddlewareContext)
	r.Post("/api/user/register", h.Register(newdb))
	r.Post("/api/user/login", h.Login(newdb))
	r.With(middleware.AuthMiddleware).Post("/api/user/orders", h.NewOrder(newdb))
	r.With(middleware.AuthMiddleware).Get("/api/user/orders", h.GetOrder(newdb))
	if err := http.ListenAndServe(flag.FlagRunAddr, r); err != nil {
		panic(err)

	}
}
