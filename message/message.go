package message

import (
	"encoding/gob"
	"io"
)

// A Message contains the actual data of a message
// sent to/from a Server, and also the type of the message.
type Message struct {
	Data interface{}
	Type string
}

// Send encodes a Message as a gob and sends it
// through a Writer.
func Send(t string, msg *Message, w io.Writer) error {
	enc := gob.NewEncoder(w)

	return enc.Encode(msg)
}

// Receive receives the last Message from a Reader.
func Receive(r io.Reader) (*Message, error) {
	var (
		dec = gob.NewDecoder(r)
		msg *Message
	)

	if err := dec.Decode(&msg); err != nil {
		return nil, err
	}

	return msg, nil
}
