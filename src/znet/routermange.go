package znet

type IRoutermanger interface {
	//返回索引对应的路由
	IndexRouter(index uint32) IRouter
	//添加路由
	AddRouter(index uint32, router IRouter)
	////删除路由
	DeleteRouter(index uint32)
}

type Routermange struct {
	//路由映射
	Routermap map[uint32]IRouter
	//负责worker取任务的消息队列

	//业务工作worker池的worker数量

}

//返回索引对应的路由,服务端通过消息绑定的ID（类型），可以执行对应不同的注册路由函数
func (R *Routermange) IndexRouter(index uint32) IRouter {
	return R.Routermap[index]
}

//添加路由
func (R *Routermange) AddRouter(index uint32, router IRouter) {
	R.Routermap[index] = router
}

////删除路由
func (R *Routermange) DeleteRouter(index uint32) {
	delete(R.Routermap, index)
}

//初始化路由
func NewRoutermange() *Routermange {
	return &Routermange{
		make(map[uint32]IRouter),
	}
}
