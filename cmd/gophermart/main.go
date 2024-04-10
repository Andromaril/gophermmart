package main

import (
	"context"
	"net/http"

	"github.com/go-chi/chi"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/andromaril/gophermmart/internal/flag"
	login "github.com/andromaril/gophermmart/internal/handler"
	storagedb "github.com/andromaril/gophermmart/internal/storage"
)

func main() {
	//var err error
	var storage storagedb.Storage
	flag.ParseFlags()
	var newdb storagedb.Storage
	//var db *sql.DB
	//if flag.Databaseflag != "" {
	db, err1 := newdb.Init(flag.Databaseflag, context.Background())
	if err1 != nil {
		panic(err1)
	}
	defer db.Close()
	//}
	//defer db.Close()
	//}
	r := chi.NewRouter()
	r.Post("/api/user/register", login.Register(storage))
	if err := http.ListenAndServe(flag.FlagRunAddr, r); err != nil {
		panic(err)

	}
}
