package utils

import (
	"encoding/json"
	"io/ioutil"
	"zinx/ziface"
)

/*

	存储一切有关 zinx 框架的全局参数，供其他模块使用
	一些参数也可以通过 zinx.json 来配置
 */
type GlobalObj struct {
	//当前Zinx的全局Server对象
	TcpServer ziface.IServer
	// 服务器的主机IP
	Host string
	// 服务器端口号
	TcpPort int
	// 服务器名称
	Name string
	// 当前 Zinx 版本号
	Version string

	// 服务器 数据包的最大值
	MaxPacketSize uint32
	//当前服务器主机允许的最大链接个数
	MaxConn int
}

var GlobalObject *GlobalObj

// Reload 读取用户配置文件
func (g *GlobalObj) Reload()  {
	data, err := ioutil.ReadFile("conf/zinx.json")
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(data, &GlobalObject)
	if err != nil {
		panic(err)
	}

}

func init()  {
	//初始化GlobalObject变量，设置一些默认值
	GlobalObject = &GlobalObj{
		Name: "Zinx ServerApp",
		Version: "v0.1",
		TcpPort: 8999,
		Host: "0.0.0.0",
		MaxConn: 1000,
		MaxPacketSize: 4096,
	}

	// 从配置文件中加载一些用户配置的参数
	GlobalObject.Reload()
}

