package znet

import "zinx/ziface"

// BaseRouter 先嵌入这几基类，然后根据需要对这个基类对方法进行重写
type BaseRouter struct {}

// 这里之所以BaseRouter的方法都为空
// 是因为有对Router不希望有PreHandle或PostHandle
// 所以Router全部继承BaseRouter的好处是，不需要实现PreHandle和PostHandle也可以实例化

// PreHandle 在处理conn业务之前的钩子方法 hook
func (br *BaseRouter) PreHandle(request ziface.IRequest) {}
// Handle 处理conn业务
func (br *BaseRouter) Handle(request ziface.IRequest) {}
// PostHandle 在处理conn业务之后的钩子方法 hook
func (br *BaseRouter) PostHandle(request ziface.IRequest){}


