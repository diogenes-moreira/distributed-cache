package distributed_cache

import (
	"bytes"
	"encoding/gob"
)

// Message is a struct that represents a message that can be sent between nodes.
type Message struct {
	CacheName string
	Key       string
	Value     interface{}
}

// ToUDP serializes the Message struct to a byte slice.
func (m *Message) ToUDP() ([]byte, error) {
	var network bytes.Buffer
	enc := gob.NewEncoder(&network)
	err := enc.Encode(m)
	if err != nil {
		return nil, err
	}
	return network.Bytes(), nil
}

// FromUDP deserializes the byte slice to a Message struct.
func (m *Message) FromUDP(data []byte) error {
	network := bytes.NewBuffer(data)
	dec := gob.NewDecoder(network)
	err := dec.Decode(m)
	if err != nil {
		return err
	}
	return nil
}

// IsCleanMessage returns true if the message is a clean message.
func (m *Message) IsCleanMessage() bool {
	return m.Key == CleanMessageKey && m.Value == nil
}

// CleanMessageKey is the key used to send a clean message.
const CleanMessageKey = "<clean>"
