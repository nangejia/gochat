package main

import (
	"net"
	"fmt"
	"os"
	"bufio"
	"strings"
	"io"
)

//循环读取服务器的消息
func ReadServer(conn net.Conn) {
	buf := make([]byte, 4096)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Println("读取服务器消息失败")
			return
		}
		fmt.Println(string(buf[:n]))
	}
}
//从键盘接收消息，并向服务器端发送
func SendMsgToServer(conn net.Conn)  {
	for{
		buf:=make([]byte,1024)
		n,_:=os.Stdin.Read(buf)
		conn.Write(buf[:n])
	}
}

func main() {
	//加载配置文件，读取服务端配置信息
	dir, error := os.Getwd()
	if error != nil {
		fmt.Println("读取本地目录错误！")
		return
	}

	fileIni, error := os.Open(dir + "/server.ini")
	if error != nil {
		fmt.Println("读取配置文件client.ini出错！")
		return
	}

	//定义服务ip信息
	var server string
	//循环从文件中读取配置信息 行的开头不是#的代表有效
	buf := bufio.NewReader(fileIni)
	for {
		b, error := buf.ReadBytes('\n')

		//读取文件中server的配置信息
		if len(b) > 8{
			userMsg:=string(b)
			if "server=" == userMsg[:7]{
				server = strings.Split(userMsg,"=")[1]
				fmt.Println(server)
			}
		}

		if error == io.EOF {
			//fmt.Println("读配置文件错误:",error)
			break
		}
	}
	fmt.Println("server:",server)
	//server:=GetServerFromIni()
	//如果读到配置
	if server != "" {
		conn, error := net.Dial("tcp", server)
		if error != nil {
			fmt.Println("连接服务器失败")
			return
		}
		defer conn.Close()

		//向服务器发送消息
		go SendMsgToServer(conn)
		//读服务器信息
		ReadServer(conn)
	} else { //没有读到配置,退出
		fmt.Println("配置文件中没有服务信息")
		return
	}
}
