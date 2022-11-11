package main

import (
	"ZakirAvrora/ChatRoom/app"
	"ZakirAvrora/ChatRoom/internals/models"
	"ZakirAvrora/ChatRoom/server"
	"flag"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	addr := flag.String("port", "8080", "Network port")
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate)
	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	room := models.NewChatRoom("general", 10, nil)
	go room.RunChatRoom()

	intSrv := server.NewServer()
	intSrv.Rooms["general"] = room

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
