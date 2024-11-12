package distributed_cache

// In this file, you can find the methods sendDelete, sendSet and sendClean
// that are used to send messages to the other nodes.
//Those methods are used internally
//in the Cache struct to send messages to the other nodes
import (
	"log"
	"net"
)

// sendDelete sends a delete message to the other nodes for a given key
func (c *Cache) sendDelete(key string) {
	message := &message{Key: key, Value: nil, CacheName: c.Name, Node: c.node}
	sendMessage(c.Address, message)
}

// sendSet sends a set message to the other nodes for a given key and value
func (c *Cache) sendSet(key string, value interface{}) {
	message := &message{Key: key, Value: value, CacheName: c.Name, Node: c.node}
	sendMessage(c.Address, message)
}

// sendClean sends a sendClean message to the other nodes
func (c *Cache) sendClean() {
	message := &message{Key: cleanMessageKey, Value: nil, CacheName: c.Name, Node: c.node}
	sendMessage(c.Address, message)
}

// sendMessage sends a message to the other nodes,
// the connection is created and closed in this function because is
// a simple UDP message
func sendMessage(address string, message *message) {
	conn := createSender(address)
	defer func(conn *net.UDPConn) {
		err := conn.Close()
		if err != nil {
			log.Println(err)
		}
	}(conn)
	data, err := message.toUDP()
	if err != nil {
		log.Println(err)
		return
	}
	_, err = conn.Write(data)
	if err != nil {
		log.Println(err)
		return
	}
}

func createSender(address string) *net.UDPConn {
	udpAddr, err := net.ResolveUDPAddr("udp", "255.255.255.255"+address)
	if err != nil {
		log.Fatal(err)
	}
	conn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		log.Fatal(err)
	}
	return conn
}
