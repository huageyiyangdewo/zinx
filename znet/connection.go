package znet

import (
	"fmt"
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

	// 该连接的处理方法 router
	Router ziface.IRouter
}

// NewConnection 创建连接的方法
func NewConnection(conn *net.TCPConn, connID uint32, router ziface.IRouter) *Connection {
	c := &Connection{
		Conn: conn,
		ConnID: connID,
		IsClosed: false,
		Router: router,
		ExitBuffChan: make(chan bool, 1),
	}
	return c
}

// StartReader 处理 conn 读取数据的Goroutine
func (c *Connection) StartReader()  {
	fmt.Println("Reader Goroutine is running...")
	defer fmt.Println(c.RemoteAddr().String(), " conn reader exit")
	defer c.Stop()

	for {
		// 读取数据到buf中
		buf := make([]byte, 512)
		_, err := c.Conn.Read(buf)
		if err != nil {
			fmt.Println("recv buf err: ", err)
			continue
		}


		// 得到当前客户端请求的 request 数据
		req := Request{
			conn: c,
			data: buf,
		}

		// 从路由 Routers 中找到 注册绑定Conn 的对应Handle
		go func(request ziface.IRequest) {
			// 模板设计模式
			c.Router.PreHandle(request)
			c.Router.Handle(request)
			c.Router.PostHandle(request)
		}(&req)

	}
}


// Start 启动连接，让当前连接工作
func (c *Connection) Start() {

	// 开启处理该链接读取到客户端数据之后的请求业务
	go c.StartReader()

	for  {
		select {
		case <- c.ExitBuffChan:
			// 得到退出消息，不在阻塞
			break
		}
	}

}

// Stop 停止连接，结束当前连接
func (c *Connection) Stop() {
	// 如果当前链接已经关闭
	if c.IsClosed == true {
		return
	}
	// todo Connection Stop() 如果用户注册了该链接的关闭回调业务，那么在此刻应该显示调用

	c.IsClosed = true
	// 关闭socket 链接
	c.Conn.Close()
	// 通知从缓冲队列读取数据的业务，该链接已经关闭
	c.ExitBuffChan <- true

	// 关闭该链接的全部管道
	close(c.ExitBuffChan)

}

// GetTCPConnection 从当前连接获取原始的socket TCPConn
func (c *Connection) GetTCPConnection() *net.TCPConn {
	return c.Conn
}

// GetConnID 获取当前连接id
func (c *Connection) GetConnID() uint32 {
	return c.ConnID
}

// RemoteAddr 获取远程客户端地址信息
func (c *Connection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}