package znet

import (
	"fmt"
	"strconv"
	"zinx/utils"
	"zinx/ziface"
)

type MsgHandler struct {
	// 存放每个MsgId 所对应的处理方法的map属性
	Apis map[uint32]ziface.IRouter
	// 业务工作worker池的数量
	WorkerPoolSize uint32
	// worker 负责去任务的消息队列
	TaskQueue []chan ziface.IRequest

}

func NewMsgHandler() *MsgHandler {
	return &MsgHandler{
		Apis: make(map[uint32]ziface.IRouter),
		WorkerPoolSize: utils.GlobalObject.WorkPoolSize,
		TaskQueue: make([]chan ziface.IRequest, utils.GlobalObject.WorkPoolSize),
	}
}

// DoMsgHandler 调用Router中具体的Handle()
func (mg *MsgHandler) DoMsgHandler(request ziface.IRequest){
	handler, ok := mg.Apis[request.GetMsgID()]
	if !ok {
		fmt.Println("api msgId = ", request.GetMsgID(), " is not FOUND!")
		return
	}

	// 执行对应的方法
	handler.PreHandle(request)
	handler.Handle(request)
	handler.PostHandle(request)

}
// AddRouter 为消息添加具体的处理逻辑
func (mg *MsgHandler) AddRouter(msgId uint32, router ziface.IRouter){
	//1 判断当前msg绑定的API处理方法是否已经存在
	if _, ok := mg.Apis[msgId]; ok {
		panic("repeated api, msgId=" + strconv.Itoa(int(msgId)))
	}
	//2 添加msg与api的绑定关系
	mg.Apis[msgId] = router
	fmt.Println("Add api msgId=", msgId)
}

// StartWorkerPool 启动 worker 工作池(只启动一次)
func (mg *MsgHandler) StartWorkerPool()  {
	// 遍历需要启动的 worker 的数量，依次启动
	for i := 0; i < int(mg.WorkerPoolSize); i++ {
		// 一个worker被启动
		// 给当前worker对应的队列开辟空间
		mg.TaskQueue[i] = make(chan ziface.IRequest, utils.GlobalObject.MaxWorkerTaskLen)
		// 启动当前Worker,阻塞的等待对应的队列是否有消息传递进来
		go mg.StartOneWorker(i, mg.TaskQueue[i])
	}
}


// StartOneWorker 为消息添加具体的处理逻辑
func (mg *MsgHandler) StartOneWorker(workerId int, taskQueue chan ziface.IRequest) {
	fmt.Println("Worker Id = ", workerId, " is started....")

	// 不断的等待队列中的消息
	for  {
		select {
		// 有消息则取出来队列中的 Request， 并执行绑定的业务方法
		case request := <- taskQueue:
			go mg.DoMsgHandler(request)
		}
	}
}

// SendMsgToTaskQueue 将消息交给TaskQueue,由worker进行处理
func (mg *MsgHandler) SendMsgToTaskQueue(request ziface.IRequest) {
	//根据ConnID来分配当前的连接应该由哪个worker负责处理
	//轮询的平均分配法则

	//得到需要处理此条连接的workerID
	workerId := request.GetConnection().GetConnID() % mg.WorkerPoolSize

	fmt.Println("Add ConnID=", request.GetConnection().GetConnID(),
		" request msgID=", request.GetMsgID(),
		"to workerID=", workerId)

	// 将请求信息发送给任务队列
	mg.TaskQueue[workerId] <- request

}