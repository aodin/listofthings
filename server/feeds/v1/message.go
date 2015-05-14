package v1

type Message interface{}

type ConnectMessage struct {
	SessionKey string `json:"session"`
}

type EventMessage struct {
	Resource string      `json:"resource"`
	Event    string      `json:"method"`
	Content  interface{} `json:"content"`
}
