package main

import (
	"ZakirAvrora/ChatRoom/app"
	"ZakirAvrora/ChatRoom/internals/models/server"
	"flag"
	"github.com/go-redis/redis"
	"log"
	"net/http"
	"os"
	"time"
)

func InitRedis() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	return rdb
}

func main() {
	addr := flag.String("port", "8080", "Network port")
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate)
	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	rdb := InitRedis()

	intSrv := server.NewServer(rdb)
	intSrv.CreateNewRoom("general", 10)

	app := &app.Application{
		ErrorLog: errorLog,
		InfoLog:  infoLog,
		Server:   intSrv,
	}

	srv := &http.Server{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  20 * time.Second,

		Addr:     ":" + *addr,
		ErrorLog: errorLog,
		Handler:  app.Routes(),
	}

	app.InfoLog.Printf("Starting server on %s", *addr)
	err := srv.ListenAndServe()
	errorLog.Fatal(err)
}
