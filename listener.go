package distributed_cache

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"time"
)

// in this file, you can find the methods createListener and handleClient
// that are used to create a connection and handle the client messages
// those methods are used internally in the Cache struct,
// to start the listener

// uDPConnInterface is an interface,
// that defines the methods of the UDPConn struct
// that are created by uncouple of the net.UDPConn struct,
// to be able to mock it for testing
type uDPConnInterface interface {
	ReadFromUDP(b []byte) (n int, addr *net.UDPAddr, err error)
	Write(b []byte) (n int, err error)
	SetReadDeadline(t time.Time) error
	Close() error
}

// startListener starts the listener to receive messages from the other nodes
func (c *Cache) startListener(ctx context.Context) {
	conn := createListener(c.Address)
	for {
		select {
		case <-ctx.Done():
			err := conn.Close()
			if err != nil {
				log.Println(err)
				return
			}
			return
		default:
			message, err := handleClient(conn)
			if err != nil {
				log.Println(err)
				continue
			}
			if message.CacheName != c.Name {
				continue
			}
			if message.isCleanMessage() {
				c.mutex.Lock()
				c.storage = make(map[string]interface{})
				c.mutex.Unlock()
				continue
			} else {
				c.Set(message.Key, message.Value)
			}
		}
		if ctx.Done() != nil {
			return
		}
	}
}

// createListener creates a connection to listen for messages
func createListener(address string) *net.UDPConn {
	udpAddr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		log.Fatal(err)
	}
	conn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return conn
}

// handleClient handles the client messages
func handleClient(conn uDPConnInterface) (*message, error) {
	buffer := make([]byte, 1024)
	err := conn.SetReadDeadline(time.Now().Add(time.Second))
	if err != nil {
		return nil, err
	}
	n, _, err := conn.ReadFromUDP(buffer)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	var message message
	network := bytes.NewBuffer(buffer[:n])
	if err := message.fromUDP(network.Bytes()); err != nil {
		log.Println(err)
		return nil, err
	}
	return &message, nil
}
