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
	if err := TemplateParseAndExecute(w, "public/index.html", nil); err != nil {
		app.ErrorLog.Println(err.Error())
		app.serverError(w, err)
	}
}
func (app *Application) wsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	client := models.NewClient("Ananymous", conn, app.Server.Rooms["general"])
	app.Server.Rooms["general"].Register <- client

	//fmt.Println("New Client joined the hub!")
	//fmt.Println(client)

	go client.WritePump()
	go client.ReadPump()

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
