package main

import (
	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/andromaril/gophermmart/internal/flag"
	"github.com/andromaril/gophermmart/internal/start"
)

func main() {
	flag.ParseFlags()
	db, newdb := start.Start()
	defer db.Close()
	go start.Update(&newdb)
	start.NewRouter(newdb)
}
