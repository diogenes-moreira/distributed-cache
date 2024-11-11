package distributed_cache

import (
	"bytes"
	"fmt"
	"log"
	"net"
	"os"
)

type UDPConnInterface interface {
	ReadFromUDP(b []byte) (n int, addr *net.UDPAddr, err error)
	Write(b []byte) (n int, err error)
	Close() error
}

func (c *Cache) StartListener() {
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
