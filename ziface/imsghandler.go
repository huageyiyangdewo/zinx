package ziface

/*
	IMsgHandler 消息管理抽象层
*/
type IMsgHandler interface {
	// DoMsgHandler 调用Router中具体的Handle()
	DoMsgHandler(request IRequest)
	// AddRouter 为消息添加具体的处理逻辑
	AddRouter(msgId uint32, router IRouter)
	// StartWorkerPool 启动worker工作池
	StartWorkerPool()
	// SendMsgToTaskQueue 将消息交给TaskQueue,由worker进行处理
	SendMsgToTaskQueue(request IRequest)
}
