package znet

import (
	"fmt"
	"io"
	"net"
	"testing"
)


// TestDataPack 只负责 测试 datapack 的拆包封包
func TestDataPack(t *testing.T) {
	// 服务端
	// 建立tcp socket 链接
	listener, err := net.Listen("tcp4", "127.0.0.1:8999")
	if err != nil {
		t.Fatal("server listen err: ", err)
	}
	// 创建服务器的 goroutine ，读取数据然后拆包
	go func() {
		// 监听 客户端请求
		conn, err := listener.Accept()
		if err != nil {
			t.Fatal("server accept err: ", err)
		}

		// 处理客户端请求
		go func(conn net.Conn) {
			// 创建 一个dp 对象
			dp := NewDataPack()
			for  {
				// 先读出msg中的head部分
				headData := make([]byte, dp.GetHeadLen())
				_, err := io.ReadFull(conn, headData) //ReadFull 会把msg填充满为止
				if err != nil {
					t.Log("read head error")
					break
				}

				// 将 headData 字节流 拆包到 msg中
				msgData, err := dp.Unpack(headData)
				if err != nil {
					t.Log("read unpack error")
					return
				}

				if msgData.GetDataLen() > 0 {
					//msg 是有data数据的，需要再次读取data数据
					msg := msgData.(*Message)
					msg.Data = make([]byte, msgData.GetDataLen())

					// 根据dataLen从 io 中读取字节流
					_, err := io.ReadFull(conn, msg.Data)
					if err != nil {
						t.Log("server unpack data error", err)
						return
					}
					fmt.Println("==> Recv Msg: ID=", msg.Id, ", len=", msg.DataLen, ", data=", string(msg.Data))
				}

			}

		}(conn)

	}()


	// 客户端
	conn, err := net.Dial("tcp4", "127.0.0.1:8999")
	if err != nil {
		t.Fatal("client dial err: ", err)
	}
	// 创建一个 封包对象 dp
	dp := NewDataPack()

	// 封装一个 msg1 包
	msg1 := Message{
		Id: 1,
		DataLen: 4,
		Data: []byte{'z', 'i', 'n', 'x'},
	}
	sendData1, err := dp.Pack(&msg1)
	if err != nil {
		t.Fatal("msg1 pack err: ", err)
	}

	// 封装一个 msg2 包
	msg2 := Message{
		Id: 1,
		DataLen: 5,
		Data: []byte{'h', 'e', 'l', 'l', 'o'},
	}
	sendData2, err := dp.Pack(&msg2)
	if err != nil {
		t.Fatal("msg2 pack err: ", err)
	}
	// 两个包拼接在一起
	sendData1 = append(sendData1, sendData2...)

	// 向服务端发送数据
	conn.Write(sendData1)
	// 阻塞
	select {

	}

}
