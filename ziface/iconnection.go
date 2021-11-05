package ziface

import "net"

// IConnection 定义连接接口
type IConnection interface {
	// Start 启动连接，让当前连接工作
	Start()

	// Stop 停止连接，结束当前连接
	Stop()

	// GetTCPConnection 从当前连接获取原始的socket TCPConn
	GetTCPConnection() *net.TCPConn

	// GetConnID 获取当前连接id
	GetConnID() uint32

	// RemoteAddr 获取远程客户端地址信息
	RemoteAddr() net.Addr

	//SendMsg 直接将Message数据发送数据给远程的TCP客户端
	SendMsg(msgId uint32, data []byte) error
}


// HandleFunc 定义统一处理连接业务的接口
// 这个是所有conn链接在处理业务的函数接口，第一参数是socket原生链接，第二个参数是客户端请求的数据，第三个参数是客户端请求的数据长度
type HandleFunc func(*net.TCPConn, []byte, int) error