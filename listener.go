package distributed_cache

import (
	"bytes"
	"fmt"
	"log"
	"net"
	"os"
)

// in this file, you can find the methods createConnection and handleClient
// that are used to create a connection and handle the client messages
// those methods are used internally in the Cache struct,
// to start the listener

// UDPConnInterface is an interface,
// that defines the methods of the UDPConn struct
// that are created by uncouple of the net.UDPConn struct,
// to be able to mock it for testing
type UDPConnInterface interface {
	ReadFromUDP(b []byte) (n int, addr *net.UDPAddr, err error)
	Write(b []byte) (n int, err error)
	Close() error
}

// startListener starts the listener to receive messages from the other nodes
func (c *Cache) startListener() {
	conn := createConnection(c.Address)
	for {
		message, err := handleClient(conn)
		if err != nil {
			log.Println(err)
			continue
		}
		if message.CacheName != c.Name {
			continue
		}
		if message.IsCleanMessage() {
			c.storage = make(map[string]interface{})
			continue
		} else {
			c.Set(message.Key, message.Value)
		}
	}
}

// createConnection creates a connection to listen for messages
func createConnection(address string) *net.UDPConn {
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
func handleClient(conn UDPConnInterface) (*Message, error) {
	buffer := make([]byte, 1024)
	n, _, err := conn.ReadFromUDP(buffer)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	var message Message
	network := bytes.NewBuffer(buffer[:n])
	if err := message.FromUDP(network.Bytes()); err != nil {
		log.Println(err)
		return nil, err
	}
	return &message, nil
}
