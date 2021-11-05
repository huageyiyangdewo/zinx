package znet

import (
	"errors"
	"fmt"
	"io"
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

	//消息管理MsgId和对应处理方法的消息管理模块
	MsgHandler ziface.IMsgHandler
}

// NewConnection 创建连接的方法
func NewConnection(conn *net.TCPConn, connID uint32, msgHandler ziface.IMsgHandler) *Connection {
	c := &Connection{
		Conn: conn,
		ConnID: connID,
		IsClosed: false,
		MsgHandler: msgHandler,
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
		// 创建拆包解包的对象
		dp := NewDataPack()
		// 读取客户端的msg head
		dataHead := make([]byte, dp.GetHeadLen())
		if _, err := io.ReadFull(c.Conn, dataHead); err != nil {
			fmt.Println("read msg head error: ", err)
			continue
		}

		//拆包，得到msgid 和 datalen 放在msg中
		msg, err := dp.Unpack(dataHead)
		if err != nil {
			fmt.Println("unpack error: ", err)
			c.ExitBuffChan <- true
			continue
		}

		//根据 dataLen 读取 data，放在msg.Data中
		var data []byte
		if msg.GetDataLen() > 0 {
			data = make([]byte, msg.GetDataLen())
			if _, err := io.ReadFull(c.Conn, data); err != nil {
				fmt.Println("read msg data error: ", err)
				c.ExitBuffChan <- true
				continue
			}
		}
		msg.SetData(data)

		// 得到当前客户端请求的 request 数据
		req := Request{
			conn: c,
			msg: msg,
		}

		//从绑定好的消息和对应的处理方法中执行对应的Handle方法
		go c.MsgHandler.DoMsgHandler(&req)

	}
}


//SendMsg 直接将Message数据发送数据给远程的TCP客户端
func (c *Connection) SendMsg(msgId uint32, data []byte) error {
	if c.IsClosed {
		return errors.New("Connection closed when send msg ")
	}

	dp := NewDataPack()
	msg, err := dp.Pack(NewMsgPackage(msgId, data))
	if err != nil {
		fmt.Println("Pack error msg id = ", msgId)
		return  errors.New("Pack error msg ")
	}


	if _, err := c.Conn.Write(msg); err != nil {
		fmt.Println("Write msg id ", msgId, " error ")
		c.ExitBuffChan <- true
		return errors.New("conn Write error")
	}

	return nil
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