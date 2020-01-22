package main

import (
	"github.com/Ledka17/TP_DB/forum/repository"
	"github.com/Ledka17/TP_DB/handler"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo"
	_ "github.com/lib/pq"
	"log"
)

const (
	PORT = "5000"
	POSTGRES = "postgres://docker:docker@localhost/docker"
)

func main() {
	e := echo.New()

	db, err := newDB()
	if err != nil {
		panic(err)
	}
	handler.NewHandler(e, repository.NewDatabaseRepository(db))

	log.Println("http server started on :5000")
	log.Fatal(e.Start(":" + PORT))
}

func newDB() (*sqlx.DB, error) {
	db, err := sqlx.Connect("pgx", POSTGRES)
	if err != nil {
		return nil, err
	}
	db.SetMaxIdleConns(7)
	db.SetMaxOpenConns(7)
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}
