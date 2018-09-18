package main

import (
	"net"
	"fmt"
	"strings"
	"time"
)

//发送消息到用户终端
func SendMsgToUser(user User, conn net.Conn) {
	for msg := range user.C {
		conn.Write([]byte(msg + "\n"))
	}
}

func MakeMsg(user User, msg string) string {
	return "[" + user.Addr + "]" + user.Name + ": " + msg
}

func HandlerConnect(conn net.Conn) {
	defer conn.Close()
	//获取在线用户的ip信息
	addr := conn.RemoteAddr().String()

	//创建一个新用户
	user := User{GetUserId(), addr, addr, make(chan string)}
	//将新用户添加到在线列表的map中
	onlineMap[addr] = &user
	//为该用户注册一个消息发送监听GO程，负责为发送消息到该用户的客户端
	go SendMsgToUser(user, conn)

	//向在线的其他用户广播新用户上线的消息
	message <- MakeMsg(user, "login")

	//用户是否退出状态
	isExit := make(chan bool)
	//是否发送消息
	hasData := make(chan bool)

	//读用户端输入的信息，并发送到全局的消息管道中
	go func() {
		buf := make([]byte, 4096)
		for {
			n, err := conn.Read(buf)
			if n == 0 {
				fmt.Println(MakeMsg(user, "退出系统"))
				isExit <- true
				return
			}
			if err != nil {
				fmt.Println("请取用户端信息错误")
				return
			}
			msg := string(buf[:n-1])
			//对用户输入的消息进行处理
			if msg == "who" && len(msg) == 3 { //查询用户在线
				for _, user := range onlineMap {
					userOnline := user.Addr + "|" + user.Name + "\n"
					conn.Write([]byte(userOnline))
				}
			} else if len(msg) > 8 && msg[:7] == "rename|" {
				name := strings.Split(msg, "|")[1]
				user.Name = name
				conn.Write([]byte("rename succee!\n"))
			} else if len(msg) > 0 { //否则为广播消息
				message <- MakeMsg(user, msg)
			}

			hasData <- true
		}
	}()

	for {
		select {
		case <-isExit:
			//删除在线用户
			delete(onlineMap, user.Addr)
			//向在线的其他用户广播用户下线消息
			message <- MakeMsg(user, "logout")
			//退出当前的GO程
			return
		case <-hasData: // 如果用户发了消息，计时重新开始
		case <-time.After(60 * time.Second):
			//删除在线用户
			delete(onlineMap, user.Addr)
			//向在线的其他用户广播用户下线消息
			message <- MakeMsg(user, " timeout !")
			//退出当前的GO程
			return
		}
	}
}

//生成当前用户的Id  由map中的元素个数决定id
func GetUserId() (id int) {
	id = len(onlineMap)
	if id == 0 {
		id = 1
	}
	return
}

//全局消息管理
func Manager() {
	for {
		//从全局消息管道中获取等待发送的消息
		msg := <-message
		//从在线列表中获取用户，将广播消息发送给每个在线用的消息管道中
		for _, user := range onlineMap {
			user.C <- msg
		}
	}
}

type User struct {
	id   int //用户编号
	Name string
	Addr string
	C    chan string //消息通道
}

//在线用户列表
var onlineMap map[string]*User = make(map[string]*User)
//全局消息通道
var message chan string = make(chan string)

func main() {
	listener, error := net.Listen("tcp", "127.0.0.1:8000")
	if error != nil {
		fmt.Println("net.Listen error.....")
		return
	}
	defer listener.Close()
	fmt.Println("服务器已启动...............")
	//启动全局消息监听
	go Manager()
	//循环监听
	for {
		conn, error := listener.Accept()
		if error != nil {
			fmt.Println("listener.Accept error.....")
			return
		}
		go HandlerConnect(conn)
	}
}
