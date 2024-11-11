package distributed_cache

import (
	"log"
	"net"
)

func (c *Cache) sendDelete(key string) {
	message := &Message{Key: key, Value: nil, CacheName: c.Name}
	sendMessage(c.Address, message)
}

func (c *Cache) sendSet(key string, value interface{}) {
	message := &Message{Key: key, Value: value, CacheName: c.Name}
	sendMessage(c.Address, message)
}

func (c *Cache) clean() {
	message := &Message{Key: CleanMessageKey, Value: nil, CacheName: c.Name}
	sendMessage(c.Address, message)
}

func sendMessage(address string, message *Message) {
	conn := createConnection(address)
	defer func(conn *net.UDPConn) {
		err := conn.Close()
		if err != nil {
			log.Println(err)
		}
	}(conn)
	data, err := message.ToUDP()
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
