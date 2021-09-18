package example

import (
	"bufio"
	"bytes"
	"fmt"
	"net"
	"time"
)

func socketClient() {
	conn, err := net.DialTimeout("tcp", "127.0.0.1:9099", time.Second*5)
	//使用带缓冲的Writer
	//缓冲默认大小为4096
	bufWriter := bufio.NewWriter(conn)

	//也可以根据自定义的缓冲大小来创建一个Writer
	//buf := bufio.NewWriterSize(conn)
	if err != nil {
		fmt.Println("connection failed to 127.0.0.1:9099 !")
		return
	}
	address := conn.LocalAddr().String()
	for {
		time.Sleep(time.Second * 2)
		//使用一个bytes缓冲来构造要发送的消息
		var bytesBuf bytes.Buffer
		bytesBuf.WriteString("from client")
		bytesBuf.WriteString(address)
		bytesBuf.WriteString(" -->")
		bytesBuf.WriteString(time.Now().String())
		_, err := bufWriter.Write(bytesBuf.Bytes())
		if err != nil {
			break
		}
		bufWriter.Flush()
	}
	conn.Close()
	fmt.Println("close connection by client !")
}
