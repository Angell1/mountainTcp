package znet

type IRouter interface {
	//Hook机制：其主要思想是提前在可能增加功能的地方埋好(预设)一个钩子，这个钩子并没有实际的意义，当我们需要重新修改或者增加这个地方的逻辑的时候，把扩展的类或者方法挂载到这个点即可。
	//在处理conn业务之前的钩子方法Hook
	PreHandle(request IRequest)
	//在处理conn业务主方法Hook
	Handle(request IRequest)
	//在处理conn业务之后的钩子方法Hook
	PostHandle(request IRequest)
}

//基类实现所有接口方法，但不是不具体写死，然后子类继承基类，重写基类方法
type BaseRouter struct{}

//在处理conn业务之前的钩子方法Hook
func (br *BaseRouter) PreHandle(request IRequest) {}

//在处理conn业务主方法Hook
func (br *BaseRouter) Handle(request IRequest) {}

//在处理conn业务之后的钩子方法Hook
func (br *BaseRouter) PostHandle(request IRequest) {}
