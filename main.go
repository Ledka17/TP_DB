package main

import (
	"fmt"
	"github.com/Ledka17/TP_DB/forum/repository"
	"github.com/Ledka17/TP_DB/handler"
	"github.com/gorilla/mux"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log"
	"net/http"
)

const (
	PORT = "5000"
	host     = "localhost"
	port     = 5432
	user     = "docker"
	password = "docker"
	dbname   = "docker"
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
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s",
		host, port, user, password, dbname)
	db, err := sqlx.Connect("pgx", psqlInfo)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}
