package app

import (
	"ZakirAvrora/ChatRoom/internals/models"
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func (app *Application) homeHandler(w http.ResponseWriter, r *http.Request) {
	if err := TemplateParseAndExecute(w, "public/home.html", nil); err != nil {
		app.ErrorLog.Println(err.Error())
		app.serverError(w, err)
	}
}

func (app *Application) roomHandler(w http.ResponseWriter, r *http.Request) {
	if err := TemplateParseAndExecute(w, "public/index.html", nil); err != nil {
		app.ErrorLog.Println(err.Error())
		app.serverError(w, err)
	}
}
func (app *Application) wsHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	nick := r.Form["nick"]
	var name string

	if len(nick) == 0 {
		name = "Ananymous"
	} else {
		name = nick[0]
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	client := models.NewClient(name, conn, app.Server.Rooms["general"])

	go client.WritePump()
	go client.ReadPump()
	app.Server.Rooms["general"].Register <- client
}

func TemplateParseAndExecute(w http.ResponseWriter, path string, data interface{}) error {
	tmpl, err := template.ParseFiles(path)
	if err != nil {
		return err
	}

	if err := tmpl.Execute(w, data); err != nil {
		return err
	}
	return nil
}
