package distributed_cache

import (
	"bytes"
	"encoding/gob"
	"testing"
)

func TestMessage_ToUDP(t *testing.T) {
	msg := &message{CacheName: "testCache", Key: "testKey", Value: "testValue"}
	data, err := msg.toUDP()
	if err != nil {
		t.Fatalf("ToUDP() error = %v", err)
	}

	var decodedMsg message
	network := bytes.NewBuffer(data)
	dec := gob.NewDecoder(network)
	err = dec.Decode(&decodedMsg)
	if err != nil {
		t.Fatalf("Decode() error = %v", err)
	}

	if decodedMsg.CacheName != msg.CacheName || decodedMsg.Key != msg.Key ||
		decodedMsg.Value != msg.Value {
		t.Errorf("Decoded message = %v, want %v", decodedMsg, msg)
	}
}

func TestMessage_FromUDP(t *testing.T) {
	msg := &message{CacheName: "testCache", Key: "testKey", Value: "testValue"}
	data, err := msg.toUDP()
	if err != nil {
		t.Fatalf("ToUDP() error = %v", err)
	}

	var decodedMsg message
	err = decodedMsg.fromUDP(data)
	if err != nil {
		t.Fatalf("FromUDP() error = %v", err)
	}

	if decodedMsg.CacheName != msg.CacheName || decodedMsg.Key != msg.Key ||
		decodedMsg.Value != msg.Value {
		t.Errorf("Decoded message = %v, want %v", decodedMsg, msg)
	}
}

func TestMessage_IsCleanMessage(t *testing.T) {
	cleanMsg := &message{Key: cleanMessageKey, Value: nil}
	if !cleanMsg.isCleanMessage() {
		t.Errorf("IsCleanMessage() = false, want true")
	}

	nonCleanMsg := &message{Key: "testKey", Value: "testValue"}
	if nonCleanMsg.isCleanMessage() {
		t.Errorf("IsCleanMessage() = true, want false")
	}
}
