package app

import (
	"ZakirAvrora/ChatRoom/internals/models/server"
	"log"
	"net/http"
)

type Application struct {
	ErrorLog *log.Logger
	InfoLog  *log.Logger
	Server   *server.Server
}

func (app *Application) Routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/", app.homeHandler)
	mux.HandleFunc("/room", app.roomHandler)
	mux.HandleFunc("/ws", app.wsHandler)
	fileServer := http.FileServer(http.Dir("./public/static"))

	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	return mux
}
