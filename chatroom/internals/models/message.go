package models

type Message struct {
	From *Client
	Msg  []byte
}
