package znet

import (
	"net"
	"zinx/ziface"
)

type Connection struct {
	// Conn 当前连接的 socket TCP 套接字
	Conn *net.TCPConn

	// ConnID 当前连接的ID，也可以作为全局唯一的SessionID，ID唯一
	ConnID uint32

	// 当前连接的关闭状态
	IsClosed bool

	// 该连接的处理方式api
	handleAPI ziface.HandleFunc
	// 告知该连接已经停止/退出 的 channel
	ExitBuffChan chan bool
}

// NewConnection 创建连接的方法
func NewConnection(conn *net.TCPConn, connID uint32, callbackApi ziface.HandleFunc) *Connection {
	c := &Connection{
		Conn: conn,
		ConnID: connID,
		IsClosed: false,
		handleAPI: callbackApi,
		ExitBuffChan: make(chan bool, 1),
	}
	return c
}
