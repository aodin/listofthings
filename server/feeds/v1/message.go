package v1

type Message interface{}

type EventMessage struct {
	Resource string      `json:"resource"`
	Event    string      `json:"method"`
	Content  interface{} `json:"content"`
}
