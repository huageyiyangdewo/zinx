package znet

import (
	"fmt"
	"net"
	"zinx/utils"
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

	// 当前Server的消息管理模块，用来绑定MsgId和对应的处理方法
	msgHandler ziface.IMsgHandler
}


// Start 开启服务
func (s *Server) Start()  {
	fmt.Printf("[START] Server name: %s,listenner at IP: %s, Port %d is starting\n", s.Name, s.IP, s.Port)
	fmt.Printf("[Zinx] Version: %s, MaxConn: %d,  MaxPacketSize: %d\n",
		utils.GlobalObject.Version,
		utils.GlobalObject.MaxConn,
		utils.GlobalObject.MaxPacketSize)

	go func() {
		// 0 启动worker工作机制
		s.msgHandler.StartWorkerPool()

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
			dealConn := NewConnection(conn, cid, s.msgHandler)
			cid++

			// 我们这里暂时做一个回显服务
			go dealConn.Start()

		}
	}()

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


func (s *Server) AddRouter(msgId uint32, router ziface.IRouter) {
	s.msgHandler.AddRouter(msgId, router)
}

// NewServer 创建一个 服务器句柄
func NewServer() ziface.IServer {
	s := &Server{
		Name: utils.GlobalObject.Name,
		IPVersion: "tcp4",
		IP: utils.GlobalObject.Host,
		Port: utils.GlobalObject.TcpPort,
		msgHandler: NewMsgHandler(),
	}

	return s
}