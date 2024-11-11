// Package distributed_cache contains all necessary files
// to add a distributed cache to your application.
// The communication between nodes using UDP.
package distributed_cache

import (
	"bytes"
	"encoding/gob"
)

// message is a struct that represents a message that can be sent between nodes.
type message struct {
	CacheName string
	Key       string
	Value     interface{}
}

// ToUDP serializes the message struct to a byte slice.
func (m *message) toUDP() ([]byte, error) {
	var network bytes.Buffer
	enc := gob.NewEncoder(&network)
	err := enc.Encode(m)
	if err != nil {
		return nil, err
	}
	return network.Bytes(), nil
}

// FromUDP deserializes the byte slice to a message struct.
func (m *message) fromUDP(data []byte) error {
	network := bytes.NewBuffer(data)
	dec := gob.NewDecoder(network)
	err := dec.Decode(m)
	if err != nil {
		return err
	}
	return nil
}

// IsCleanMessage returns true if the message is a clean message.
func (m *message) isCleanMessage() bool {
	return m.Key == cleanMessageKey && m.Value == nil
}

// cleanMessageKey is the key used to send a clean message.
const cleanMessageKey = "<clean>"
