package apps

import (
	"log"
	"net"
	"time"
)

type Message struct {
	Name    string
	Content string
	RoomNum int
}

type ClientInfo struct {
	Name       string
	CreateTime time.Time
	conn       net.Conn
	RoomNum    int
	Destroy    bool
}

func (u *ClientInfo) Close() {
	if u == nil {
		return
	}
	delete(chatServerMgr.chatServer[u.RoomNum], u.Name)
	u.conn.Close()
	u.Destroy = true
}

func (u *ClientInfo) Write(b []byte) {
	_, err := u.conn.Write(b)
	if err != nil {
		log.Printf("broad message failed: %v\n", err)
		u.Close()
	}
}
