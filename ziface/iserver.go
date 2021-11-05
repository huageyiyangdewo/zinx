package ziface


// IServer 定义服务器接口
type IServer interface {
	// Start 启动服务器
	Start()
	// Stop 停止服务器
	Stop()
	// Serve 开启业务服务器
	Serve()

	// AddRouter 路由功能：给当前服务注册一个路由业务方法，供客户端处理使用
	AddRouter(msgId uint32, router IRouter)
}