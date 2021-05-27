package main

import (
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	Run()
}

func Run() {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", ":1526")
	if err != nil {
		log.Printf("Resolve tcp addr failed: %v\n", err)
		return
	}

	// 向服务器拨号
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		log.Printf("Dial to server failed: %v\n", err)
		return
	}

	// 设置Name和进入的聊天室序号
	var nickName string
	var roomNum int
	fmt.Println("请输入昵称 房间号:")
	fmt.Scanln(&nickName, &roomNum)
	input := fmt.Sprintf("%v %v", nickName, roomNum)
	_, err = conn.Write([]byte(input))
	if err != nil {
		conn.Close()
		return
	}

	// 向服务器发消息
	go SendMsg(conn)

	// 接收来自服务器端的消息
	buf := make([]byte, 1024)
	for {
		length, err := conn.Read(buf)
		if err != nil {
			log.Printf("recv server msg failed: %v\n", err)
			conn.Close()
			os.Exit(0)
			break
		}

		fmt.Println(string(buf[0:length]))
	}
}

// 向服务器端发消息
func SendMsg(conn net.Conn) {
	for {
		var input string
		fmt.Scanln(&input)

		if input == "/q" || input == "/quit" {
			fmt.Println("Byebye ...")
			conn.Close()
			os.Exit(0)
		}

		// 只处理有内容的消息
		if len(input) > 0 {
			_, err := conn.Write([]byte(input))
			if err != nil {
				conn.Close()
				break
			}
		}
	}
}
