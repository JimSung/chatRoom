package apps

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
	"time"
)

type ChatServerMgr struct {
	chatServer map[int]map[string]*ClientInfo
}

// 启动服务器
func Run(roomNux int) {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", ":1526")
	if err != nil {
		log.Printf("resolve tcp addr failed: %v\n", err)
		return
	}

	chatServerMgr = &ChatServerMgr{
		chatServer: make(map[int]map[string]*ClientInfo),
	}

	for i := 1; i <= roomNux; i++ {
		chatServerMgr.chatServer[i] = make(map[string]*ClientInfo)
	}

	// 监听
	listener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		log.Printf("listen tcp port failed: %v\n", err)
		return
	}

	// 消息通道
	messageChan := make(chan []byte, 10)

	// 广播消息
	go BroadMessages(messageChan)

	// 启动
	for {
		conn, err := listener.AcceptTCP()
		if err != nil {
			log.Printf("Accept failed:%v\n", err)
			continue
		}

		user := &ClientInfo{
			CreateTime: time.Now(),
			conn:       conn,
		}

		for {
			buf := make([]byte, 1024)
			length, err := conn.Read(buf)
			if err != nil {
				log.Printf("read client message failed:%v\n", err.Error())
				conn.Close()
				break
			}
			msg := string(buf[0:length])
			paraList := strings.Split(msg, " ")
			if len(paraList) != 2 {
				conn.Close()
				break
			}
			user.Name = paraList[0]
			user.RoomNum, err = strconv.Atoi(paraList[1])
			if err != nil {
				conn.Close()
				break
			}
			errMsg := ""
			if _, ok := chatServerMgr.chatServer[user.RoomNum]; !ok {
				errMsg = "房间不存在"
			}
			if _, ok := chatServerMgr.chatServer[user.RoomNum][user.Name]; ok {
				errMsg = "该昵称已存在"
			}
			if errMsg != "" {
				user.Write([]byte(errMsg))
				user.Close()
			}
			break
		}

		SendHistory(user)
		if user.Destroy {
			break
		}

		chatServerMgr.chatServer[user.RoomNum][user.Name] = user
		// 处理消息
		go Handler(user, messageChan)
	}
}

// 发广播
func BroadMessages(messages chan []byte) {
	for {

		msg := <-messages
		fmt.Println(string(msg))

		m := &Message{}
		er := json.Unmarshal(msg, m)
		if er != nil {
			log.Printf("Unmarshal field:%v\n", er.Error())
			break
		}
		content := fmt.Sprintf("%v:%v", m.Name, m.Content)

		for _, user := range chatServerMgr.chatServer[m.RoomNum] {
			if user.Name == m.Name {
				continue
			}
			user.Write([]byte(content))
		}
	}
}

// 处理客户端发到服务端的消息，将其扔到通道中
func Handler(u *ClientInfo, messages chan []byte) {
	fmt.Println("connect from client ", u.Name)

	buf := make([]byte, 1024)
	for {
		length, err := u.conn.Read(buf)
		if err != nil {
			log.Printf("read client message failed:%v\n", err.Error())
			u.Close()
			break
		}
		if length == 0 {
			continue
		}

		message := &Message{
			Name:    u.Name,
			Content: string(buf[0:length]),
			RoomNum: u.RoomNum,
		}

		// 检测gm命令
		isCmd := false
		if _, ok := cmdMgr.CmdHandler[message.Content]; ok {
			handler := cmdMgr.CmdHandler[message.Content]
			result := handler(u)
			if result != "" {
				isCmd = true
				message.Content = result
			}
		}

		// 过滤脏词
		if !isCmd {
			s := filterMgr.Filter(message.Content)
			message.Content = s
		}

		if isCmd {
			u.Write([]byte(message.Content))
		} else {
			chatMgr.Add(message)
			b, er := json.Marshal(message)
			if er != nil {
				log.Printf("Marshal field:%v\n", er.Error())
				break
			}
			messages <- b
		}
	}
}

func SendHistory(u *ClientInfo) {
	for _, msg := range chatMgr.GetHistory() {
		content := fmt.Sprintf("%v:%v\n", msg.Name, msg.Content)
		u.Write([]byte(content))
	}
}
