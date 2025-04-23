package main

import (
	"net"
	"strings"
)

type User struct {
	Name   string
	Addr   string
	Ch     chan string
	conn   net.Conn
	server *Server
}

// 监听用户通道
func (u *User) Listen() {
	for {
		lis := <-u.Ch
		u.conn.Write([]byte(lis))
	}
}

// 用户上线
func (u *User) Online() {
	u.server.mapLock.Lock() //上锁
	u.server.Online[u.Name] = u
	u.server.mapLock.Unlock() //解锁
	//广播用户上线消息
	u.server.Broadcast(u, "已上线")
}

// 用户下线
func (u *User) Offline() {
	u.server.mapLock.Lock() //上锁
	delete(u.server.Online, u.Name)
	u.server.mapLock.Unlock() //解锁
	//广播用户上线消息
	u.server.Broadcast(u, "下线")
}

// 给对应的user发送消息
func (u *User) Senduser(msg string) {
	u.conn.Write([]byte(msg))
}

// 用户处理消息
func (u *User) Domsg(msg string) {
	if msg == "who" {
		u.server.mapLock.Lock()
		for _, U := range u.server.Online {
			if U == u {
				onlinemsg := "[" + U.Addr + "](" + U.Name + ")本人" + "[在线...]" + "\r\n"
				u.Senduser(onlinemsg)
			} else {
				onlinemsg := "[" + U.Addr + "](" + U.Name + ")" + "[在线...]" + "\r\n"
				u.Senduser(onlinemsg)
			}

		}
		u.server.mapLock.Unlock()
	} else if len(msg) > 7 && msg[:7] == "rename|" { //判断是否是改名 格式rename|名字
		newName := strings.Split(msg, "|")[1]
		_, ok := u.server.Online[newName]
		if ok {
			u.Senduser("用户名重复，请重试！\r\n")
		} else {
			u.server.mapLock.Lock()
			delete(u.server.Online, u.Name)
			u.server.Online[newName] = u
			u.server.mapLock.Unlock()
			u.Name = newName
			u.Senduser("您的用户名已更新为[" + newName + "]\r\n")
		}

	} else if len(msg) > 3 && msg[:3] == "to|" { //判断是否是私发 格式to|名字|消息
		parts := strings.Split(msg, "|")
		if len(parts) == 3 {
			name := parts[1]
			if name == "" {
				u.Senduser("发送格式错误\r\n")
				return
			}
			touser, ok := u.server.Online[name]
			if !ok {
				u.Senduser("用户不存在\r\n")
				return
			}
			tomsg := strings.Split(msg, "|")[2]
			if tomsg == "" {
				u.Senduser("消息为空，请重复\r\n")
				return
			}
			touser.Senduser("[" + u.Name + "]给你发消息：" + tomsg + "\r\n")
			u.Senduser("消息已发送\r\n")
		} else {
			u.Senduser("格式错误，应该为 to|名字|消息 \r\n")
		}

	} else {
		u.server.Broadcast(u, msg)
	}
}

// 创建一个用户
func NewUser(conn net.Conn, s *Server) *User {
	userAddr := conn.RemoteAddr().String()
	usr := &User{userAddr, userAddr, make(chan string), conn, s}
	go usr.Listen()
	return usr
}
