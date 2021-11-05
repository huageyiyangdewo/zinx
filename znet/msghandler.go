package znet

import (
	"fmt"
	"strconv"
	"zinx/ziface"
)

type MsgHandler struct {
	// 存放每个MsgId 所对应的处理方法的map属性
	Apis map[uint32]ziface.IRouter
}

func NewMsgHandler() *MsgHandler {
	return &MsgHandler{
		Apis: make(map[uint32]ziface.IRouter),
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
