package main

import (
	"TP_DB/forum/repository"
	"TP_DB/handler"
	"github.com/gorilla/mux"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"os"
)

var (
	PORT = "5000"
)

func main() {
	db, err := newDB()
	if err != nil {
		panic(err)
	}
	repository.NewDatabaseRepository(db)
	r := mux.NewRouter()
	handler.NewHandler(r, repository.NewDatabaseRepository(db))

	log.Fatal(http.ListenAndServe(":" + PORT, r))
}

func newDB() (*sqlx.DB, error) {
	db, err := sqlx.Connect("pgx", os.Getenv("POSTGRES_DSN"))
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}
