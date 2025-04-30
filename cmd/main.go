package main

import (
	"github.com/AndreySirin/avito-backend-assignment-2023/internal/server"
	"github.com/AndreySirin/avito-backend-assignment-2023/internal/storage"
)

const (
	USER     = "admin"
	PASSWORD = "secret"
	DB       = "mydb"
	port     = "5432"
	host     = "localhost"
)

func main() {
	lg := logg.NewLogg()
	lg.Info("Start server")

	db, err := storage.NewStorage(lg, USER, PASSWORD, DB, host, port)
	if err != nil {
		lg.Error("ошибка при подключении к базе", err)
	}
	defer db.Close()

	srv := server.NewServer(lg, ":8080", db)
	err = srv.Start()
	if err != nil {
		lg.Error("ошибка при старте", err)
	}
}
