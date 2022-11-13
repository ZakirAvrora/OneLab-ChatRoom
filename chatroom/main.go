package main

import (
	"ZakirAvrora/ChatRoom/app"
	"ZakirAvrora/ChatRoom/internals/models/server"
	"flag"
	"fmt"
	"github.com/go-redis/redis"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

func InitRedis(address string) *redis.Client {
	address = strings.TrimSpace(address)
	if address == "" {
		address = "localhost"
	}
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf(address + ":6379"),
		Password: "",
		DB:       0,
	})
	return rdb
}

func main() {
	addr := flag.String("port", "8080", "Network port")
	redisAddress := flag.String("redis-address", "localhost", "Redis address")
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate)
	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	rdb := InitRedis(*redisAddress) // change on localhost
	infoLog.Printf("using %s redis address", *redisAddress)

	intSrv := server.NewServer(rdb)
	intSrv.CreateNewRoom("general", 10)

	application := &app.Application{
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
		Handler:  application.Routes(),
	}

	application.InfoLog.Printf("Starting server on %s", *addr)
	err := srv.ListenAndServe()
	errorLog.Fatal(err)
}
