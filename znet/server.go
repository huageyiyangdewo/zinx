package znet

import (
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

	// 3 启动Server网络连接业务
	for  {
		// 3.1 阻塞等待客服端连接
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Accept err: ", err)
			continue
		}

		// 3.2 todo Server.Start() 设置服务器最大连接控制，如果超过最大连接，那么关闭此连接
		// 3.3 todo Server.Start() 处理该新连接请求的 业务 方法，此时应该有 handler 和 conn 是绑定的

		// 我们这里暂时做一个最大512 字节的回显服务
		for  {
			buf := make([]byte, 512)
			cnt, err := conn.Read(buf)
			if err != nil {
				fmt.Println("recv buf err ", err)
				continue
			}

			if _, err := conn.Write(buf[:cnt]); err != nil {
				fmt.Println("write back buf err ", err)
				continue
			}
		}

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

// NewServer 创建一个 服务器句柄
func NewServer(Name string) ziface.IServer {
	s := &Server{
		Name: Name,
		IPVersion: "tcp4",
		IP: "0.0.0.0",
		Port: 8999,
	}

	return s
}