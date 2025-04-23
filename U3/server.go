package main

import (
	"bufio"
	. "fmt"
	"io"
	"net"
	"strings"
	"sync"
	"time"
)

// 服务端结构体
type Server struct {
	Ip      string
	Port    int
	Online  map[string]*User //用户在线列表
	mapLock sync.RWMutex
	Message chan string //广播通道
}

// 创建server接口
func NewServer(ip string, port int) *Server {
	return &Server{Ip: ip, Port: port, Online: make(map[string]*User), Message: make(chan string)}
}

// 广播消息
func (s *Server) Broadcast(user *User, msg string) {
	sendmsg := "[" + user.Addr + "](" + user.Name + "):" + msg + "\r\n"
	s.Message <- sendmsg
}

// 监听广播消息，监听到后发送给所有的用户
func (s *Server) Lismsg() {
	for {
		msg := <-s.Message
		s.mapLock.Lock()
		for _, u := range s.Online {
			u.Ch <- msg
		}
		s.mapLock.Unlock()
	}
}

func (s *Server) Handle(conn net.Conn) {
	user := NewUser(conn, s)
	user.Online()
	//监听用户是否活跃
	isline := make(chan bool)
	//接收客户端发送的消息
	go func() {
		reader := bufio.NewReader(conn) // 创建缓冲读取器
		for {
			msg, err := reader.ReadString('\n') // 按行读取，直到 \n
			if err != nil {
				if err == io.EOF {
					user.Offline()
				}
				return
			}
			Printf(msg)
			// 去除末尾的 \r\n
			msg = strings.TrimSuffix(msg, "\r\n")
			user.Domsg(msg)
			isline <- true
		}
	}()
	for {
		//防止函数退出
		select {
		case <-isline:
		case <-time.After(time.Second * 5):
			user.Senduser("你被踢下线了，请退出后重启")
			user.Offline()
			close(user.Ch)
			conn.Close()
			return
		}
	}
}

// 启动服务器
func (s *Server) Start() {
	Println("服务器已启动")
	//监听
	listener, err := net.Listen("tcp", Sprintf("%s:%d", s.Ip, s.Port))
	if err != nil {
		Println("Error listening:", err.Error())
		return
	}
	defer listener.Close()
	go s.Lismsg()
	//启动监听
	for {
		conn, err := listener.Accept()
		if err != nil {
			Println("Error accepting: ", err.Error())
			continue
		}
		go s.Handle(conn)
	}

}
