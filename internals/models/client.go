package models

import (
	"github.com/gorilla/websocket"
	"log"
	"time"
)

type Client struct {
	Nick string
	Room *ChatRoom
	Conn *websocket.Conn
	Msg  chan Message
}

func NewClient(nick string, conn *websocket.Conn, room *ChatRoom) *Client {
	return &Client{
		Nick: nick,
		Conn: conn,
		Room: room,
		Msg:  make(chan Message),
	}
}

const (
	writeWait      = 5 * time.Second
	pongWait       = 3 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 10000
)

func (c *Client) ReadPump() {
	defer func() {
		c.Room.Unregister <- c
		close(c.Msg)
		c.Conn.Close()
		c.Room.Broadcast <- Message{From: c, Msg: []byte(MsgUserLeft(c))}
	}()

	c.Conn.SetReadLimit(maxMessageSize)
	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(string) error { c.Conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	for {
		_, jsonMessage, err := c.Conn.ReadMessage()

		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("unexpected close error: %v", err)
			}
			break
		}
		c.Room.Broadcast <- Message{From: c, Msg: []byte(MsgTimeUser(c) + string(jsonMessage))}
	}
}

func (c *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)

	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.Msg:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message.Msg)

			n := len(c.Msg)
			for i := 0; i < n; i++ {
				w.Write([]byte{'\n'})
				msg := <-(c.Msg)
				w.Write(msg.Msg)
			}

			if err := w.Close(); err != nil {
				return
			}

		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func MsgTimeUser(c *Client) string {
	return "[" + time.Now().Format("2006-01-02 15:04:05") + "]" + "[" + c.Nick + "]: "
}

func MsgUserLeft(c *Client) string {
	return c.Nick + " is left our chatroom..."
}

func MsgUserIn(c *Client) string {
	return c.Nick + " is entered our chatroom..."
}
