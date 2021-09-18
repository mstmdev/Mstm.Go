package example

import (
	"fmt"
	"net"
	"time"
)

func socketServer() {
	listener, err := net.Listen("tcp", "127.0.0.1:9099")
	if err != nil {
		fmt.Println("create listener faile !")
		return
	}
	for {
		conn, err := listener.Accept()

		if err != nil {
			fmt.Println("client connect faied !")
		} else {
			fmt.Println("a  client connected by ", conn.RemoteAddr().Network(), "  -->", conn.RemoteAddr().String())
		}

		go readMsg(conn)
	}
}

func readMsg(conn net.Conn) {
	receiveData := make([]byte, 100)
	for {
		//设置超时时间
		err := conn.SetDeadline(time.Now().Add(time.Second * 3))
		count, err := conn.Read(receiveData)
		if err != nil {
			fmt.Println(err.Error())
			conn.Close()
			return
		}
		if count <= 0 {
			fmt.Println("read end")
			return
		} else {
			fmt.Println(string(receiveData))

		}
	}
}
