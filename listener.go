package distributed_cache

import (
	"bytes"
	"context"
	"fmt"
	"github.com/google/uuid"
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

type iCache interface {
	set(key string, value interface{})
	clean()
	getAddress() string
	getName() string
	getNode() uuid.UUID
}

// startListener starts the listener to receive messages from the other nodes
func startListener(c iCache, ctx context.Context) {
	conn := createListener(c.getAddress())
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
			if message.CacheName != c.getName() || message.Node == c.getNode() {
				continue
			}
			if message.isCleanMessage() {
				c.clean()
				continue
			} else {
				c.set(message.Key, message.Value)
			}
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
