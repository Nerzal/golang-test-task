package domain

import "encoding/json"

// Message is used to publish a message on the rabbitMQ.
type Message struct {
	Sender  string `json:"sender" binding:"required"`
	Message string `json:"message" binding:"required"`
}

func (m Message) MarshalBinary() (data []byte, err error) {
	bytes, err := json.Marshal(m)
	return bytes, err
}
