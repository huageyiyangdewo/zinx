package znet

import (
	"errors"
	"fmt"
	"net"
	"zinx/ziface"
)

// Server iServer 接口实现，定义一个接口实现类
type Server struct {
	// 服务器名称
	Name string
	// 服务器IP版本
	IPVersion string
	// 服务器IP
	IP string
	// 服务器端口
	Port int

	// 当前Server由用户绑定的回调 router,也就是 Server 注册的链接对应的处理业务
	Router ziface.IRouter
}

// CallbackToClient 定义当前客户端链接的handle api
func CallbackToClient(conn *net.TCPConn, data []byte, cnt int) error {
	// 回显业务
	fmt.Println("[Conn handle] CallbackToClient ...")
	if _, err := conn.Write(data[:cnt]); err != nil {
		fmt.Println("write back buf err: ", err)
		return errors.New("CallbackToClient error")
	}

	return nil
}


// Start 开启服务
func (s *Server) Start()  {
	fmt.Printf("[Start] Server listener at IP:%s, Port:%s, is starting \n", s.IP, s.Port)

	// 1 获取一个tcp的addr
	addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
	if err != nil {
		fmt.Println("resolve tcp addr err: ", err)
		return
	}

	// 2 监听服务器地址
	listener, err := net.ListenTCP(s.IPVersion, addr)
	if err != nil {
		fmt.Println("listen: ", s.IPVersion, "err: ", err)
	}
	fmt.Println("start Zinx server ", s.Name, "succ, now listening...")
	var cid uint32
	cid = 0

	// 3 启动Server网络连接业务
	for  {
		// 3.1 阻塞等待客服端连接
		conn, err := listener.AcceptTCP()
		if err != nil {
			fmt.Println("Accept err: ", err)
			continue
		}

		// 3.2 todo Server.Start() 设置服务器最大连接控制，如果超过最大连接，那么关闭此连接

		// 3.3 处理该新连接请求的 业务 方法，此时应该有 handler 和 conn 是绑定的
		dealConn := NewConnection(conn, cid, CallbackToClient)
		cid++

		// 我们这里暂时做一个最大512 字节的回显服务
		go dealConn.Start()

	}
}

// Stop 停止服务
func (s *Server) Stop()  {

}

// Serve 运行服务
func (s *Server) Serve()  {
	s.Start()

	// todo Server.Serve() 是否在启动的时候做其他的处理，可以在这里添加

	select {}

}


func (s *Server) AddRouter(router ziface.IRouter) {
	s.Router = router
}

// NewServer 创建一个 服务器句柄
func NewServer(Name string) ziface.IServer {
	s := &Server{
		Name: Name,
		IPVersion: "tcp4",
		IP: "0.0.0.0",
		Port: 8999,
		Router: nil,
	}

	return s
}