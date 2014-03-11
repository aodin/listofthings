package server

import ()

type Message struct {
	Id      int64       `json:"id"`
	Body    string      `json:"body"`
	Content interface{} `json:"content"`
}

type ReceivedMessage struct {
	Method  string `json:"method"`
	Content string `json:"content"`
}

type ThingMessage struct {
	Method string `json:"method"`
	Item   *Thing `json:"content"`
}
