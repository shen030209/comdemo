package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"
)

var serverip string
var serverport int

type Client struct {
	Serverip   string
	Serverport int
	Name       string
	conn       net.Conn
	Flag       int
}

// 链接服务器
func NewClient(ip string, port int) *Client {
	client := &Client{ip, port, "", nil, -1}
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", ip, port))
	if err != nil {
		fmt.Println("链接服务器失败", err)
		return nil
	}
	client.conn = conn
	return client
}

// 选择功能
func (c *Client) menu() bool {
	fmt.Println("c.Flag:", c.Flag)
	var flag int
	fmt.Println("1.公聊")
	fmt.Println("2.私聊")
	fmt.Println("3.更改名字")
	fmt.Println("0.退出")
	// 正确读取输入并处理错误
	_, err := fmt.Scanln(&flag)
	if err != nil {
		fmt.Println("输入无效，请重新输入")
		return false
	}
	// 检查输入范围
	if flag >= 0 && flag <= 3 && c.Flag != 0 {
		c.Flag = flag
		return true
	} else if c.Flag == 0 {
		return true
	} else {
		fmt.Println("没有这个选项")
		return false
	}

}
func (c *Client) Run() {
	for c.Flag != 0 {
		for c.menu() != true {
		}
		switch c.Flag {
		case 1:
			fmt.Println("<公聊>")
			for c.Puchat() == true {
			}
			break
		case 2:
			fmt.Println("<私聊>")
			for c.Prchat() == true {
			}
			break
		case 3:
			fmt.Println("<更改名字>")
			c.Rename()
			break
		case 0:
			c.conn.Close()
			fmt.Println("已退出")
			return
		}

	}
}

// 监听服务器消息
func (c *Client) Lisserver() {
	//监听
	io.Copy(os.Stdout, c.conn)

}

// 改名
func (c *Client) Rename() {
	var name string
	fmt.Println("输入用户名：")
	for name == "" {
		fmt.Scanln(&name)
		if name == "" {
			fmt.Println("用户名不能为空,请重新输入（输入quit退出）")
		}
		if name == "quit" {
			return
		}
	}
	c.Name = name
	msg := "rename|" + c.Name + "\r\n"
	_, err := c.conn.Write([]byte(msg))
	if err != nil {
		fmt.Println("发送错误", err)
		return
	} else {
		return
	}
}

// 公聊
func (c *Client) Puchat() bool {
	var msg string
	fmt.Println("输入内容(输入quit退出)：")
	fmt.Scanln(&msg)
	if msg == "quit" {
		return false
	} else {
		msg := msg + "\r\n"
		_, err := c.conn.Write([]byte(msg))
		if err != nil {
			fmt.Println("发送错误", err)
			return false
		} else {
			return true
		}
	}

}

// 私聊
func (c *Client) Prchat() bool {
	var name string
	var msg string
	fmt.Println("输入以下私聊人的名字(输入quit退出)：")
	c.conn.Write([]byte("who" + "\r\n"))
	fmt.Scanln(&name)
	if name == "quit" {
		return false
	} else {
		fmt.Println("输入私聊内容(输入quit退出)：")
		fmt.Scanln(&msg)
		str := "to|" + name + "|" + msg + "\r\n"
		_, err := c.conn.Write([]byte(str))
		if err != nil {
			fmt.Println("发送错误", err)
			return false
		} else {
			return true
		}
	}

}

// 自定义：输入go run ./client -ip "自定义ip" -port "自定义port"
// 默认：输入go run ./client
func init() {
	flag.StringVar(&serverip, "ip", "127.0.0.1", "(默认为127.0.0.1)")
	flag.IntVar(&serverport, "port", 8888, "默认为8888")
}
func main() {
	flag.Parse()
	client := NewClient(serverip, serverport)
	if client == nil {
		return
	}
	go client.Lisserver()
	fmt.Println("服务器链接成功")
	client.Run()
}
