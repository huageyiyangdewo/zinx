package znet

type Message struct {
	Id      uint32 // 消息的id
	DataLen uint32 // 消息的长度
	Data    []byte // 消息的内容
}


// GetMsgId 获取消息ID
func (m *Message) GetMsgId() uint32 {
	return m.Id
}
// GetDataLen 获取消息数据段长度
func (m *Message) GetDataLen() uint32{
	return m.DataLen
}
// GetData 获取消息
func (m *Message) GetData() []byte {
	return m.Data
}

// SetMsgId 设置消息ID
func (m *Message) SetMsgId(id uint32){
	m.Id = id
}
// SetData 设置消息
func (m *Message) SetData(data []byte){
	m.Data = data
}
// SetDateLen 设置消息长度
func (m *Message) SetDateLen(len uint32){
	m.DataLen = len
}