package main

import (
	"os"
	"fmt"
	"bufio"
	"strings"
	"io"
)

func GetServerFromIni() (server string) {
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
	defer fileIni.Close()
	//定义服务ip信息
	//var server string
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

	return
}
