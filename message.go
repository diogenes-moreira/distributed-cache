package distributed_cache

import (
	"bytes"
	"encoding/gob"
)

// Message is a struct that represents a message that can be sent to the cache.
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

func (m *Message) FromUDP(data []byte) error {
	network := bytes.NewBuffer(data)
	dec := gob.NewDecoder(network)
	err := dec.Decode(m)
	if err != nil {
		return err
	}
	return nil
}

func (m *Message) IsCleanMessage() bool {
	return m.Key == CleanMessageKey && m.Value == nil
}

const CleanMessageKey = "<clean>"
