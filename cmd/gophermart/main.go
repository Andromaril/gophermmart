package main

import (
	"context"
	"database/sql"
	"net/http"

	"github.com/go-chi/chi"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/andromaril/gophermmart/internal/flag"
	login "github.com/andromaril/gophermmart/internal/handler"
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
	r.Post("/api/user/register", login.Register(newdb))
	if err := http.ListenAndServe(flag.FlagRunAddr, r); err != nil {
		panic(err)

	}
}
