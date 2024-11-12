package distributed_cache

import (
	"net"
	"testing"
	"time"
)

type MockUDPConn struct {
	data []byte
}

func (m *MockUDPConn) SetReadDeadline(time time.Time) error {
	return nil
}

func (m *MockUDPConn) ReadFromUDP(b []byte) (n int, addr *net.UDPAddr, err error) {
	copy(b, m.data)
	return len(m.data), &net.UDPAddr{}, nil
}

func (m *MockUDPConn) Write(b []byte) (n int, err error) {
	m.data = b
	return len(b), nil
}

func (m *MockUDPConn) Close() error {
	return nil
}

func createMockConnection(data []byte) *MockUDPConn {
	return &MockUDPConn{data: data}
}

func TestHandleClient(t *testing.T) {
	message := &message{CacheName: "testCache", Key: "testKey", Value: "testValue"}
	data, err := message.toUDP()
	if err != nil {
		t.Fatalf("ToUDP() error = %v", err)
	}

	mockConn := createMockConnection(data)
	receivedMessage, err := handleClient(mockConn)
	if err != nil {
		t.Fatalf("handleClient() error = %v", err)
	}

	if receivedMessage.CacheName != message.CacheName || receivedMessage.Key != message.Key ||
		receivedMessage.Value != message.Value {
		t.Errorf("Expected message %v, got %v", message, receivedMessage)
	}
}
