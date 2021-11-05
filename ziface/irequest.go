package ziface

/*
	IRequest：
	实际上是把请求的链接信息和请求的数据 包装到了 Request 里
 */
type IRequest interface {
	// GetConnection 获取请求的链接信息
	GetConnection()  IConnection

	// GetData 获取请求的数据
	GetData()  []byte

	// GetMsgID 获取请求的消息ID
	GetMsgID() uint32
}
