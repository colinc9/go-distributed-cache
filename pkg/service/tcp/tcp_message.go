package tcp

import (
	"encoding/gob"
	"io"
)

type MessageType int64

const (
	Test MessageType = iota
	Get
	Set
)

type Message struct {
	Type MessageType
	Key  interface{}
	Value  interface{}
}

func Decoder (r io.Reader) (*Message, error) {
	var msg *Message
	err := gob.NewDecoder(r).Decode(&msg)
	return msg, err
}

func Encode (w io.Writer, msg *Message)  error {
	return gob.NewEncoder(w).Encode(&msg)
}