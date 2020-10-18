package znet

import (
	"fmt"
	"net"
)

type IServer interface {
	//启动
	Start()
	//停止
	Stop()
	//运行
	Run()
	//路由功能：给当前服务注册一个路由方法，供客户端链接处理使用
	AddRouter(index uint32, router IRouter)
	//得到路由管理器
	GetRoutermanger() IRoutermanger
	//得到链接管理器
	GetConnManager() IConnManager
}

//IServer接口的实现，定义一个Server的服务模块
type Server struct {
	//服务器的名称
	Name string
	//IP版本
	IPversion string
	//服务器的IP
	IP string
	//服务器的端口
	Port string
	//最大链接数
	MaxConn int64
	//路由管理器
	Routermange IRoutermanger
	//链接管理器
	ConnManager IConnManager
	//工作池
	Pool *WorkerPool
}

//启动服务器
func (s *Server) Start() {
	s.Pool.StartWorkPool()
	//在这里使用go协程让服务端阻塞接收客户端的消息为异步。
	//1.获取一个tcp的addr
	addr, err := net.ResolveTCPAddr(s.IPversion, fmt.Sprintf("%s:%s", s.IP, s.Port))
	if err != nil {
		fmt.Println("resolve tcp addr error:", err)
		return
	}
	//2.监听服务器的地址
	listenner, err := net.ListenTCP(s.IPversion, addr)
	if err != nil {
		fmt.Printf("listen:%s err:%s", s.IPversion, err)
		return
	}
	fmt.Println("Start Zinx serversucc:", s.Name)

	var cid uint32
	cid = 0
	go func() {
		//3.阻塞的等待客户端链接，处理客户端链接业务
		for {
			//如果有客户端链接过来，阻塞会返回
			conn, err := listenner.AcceptTCP()
			if err != nil {
				fmt.Println("Accept err", err)
				continue
			}
			//判断是否超过最大链接数
			if int64(s.ConnManager.GetConnLen()) >= s.MaxConn {
				err := conn.Close()
				if err != nil {
					fmt.Println("拒绝用户失败！")
				}
			} else {
				//将每个客户端链接与服务端路由模块进行绑定
				//c:封装每个客户端链接信息，绑定链接、id、及服务器的路由管理器
				Conn := NewConnection(s, conn, cid, s.Routermange)
				cid++
				//启动链接进行业务处理
				s.Pool.Put(Conn)
			}
		}
	}()

}

//停止服务器
func (s *Server) Stop() {
	s.ConnManager.ClearConn()
	//将服务器的资源、状态或者已经开辟的链接信息回收

}

//运行服务器
func (s *Server) Run() {
	//启动Server的服务功能
	s.Start()

	//做一些启动服务器之后的业务//

	//阻塞状态
	select {}
}

//添加路由功能
func (s *Server) AddRouter(index uint32, router IRouter) {
	s.Routermange.AddRouter(index, router)
	fmt.Println("Add Router Succ!!!")
}

func (s *Server) GetRoutermanger() IRoutermanger {
	return s.Routermange
}

func (s *Server) GetConnManager() IConnManager {
	return s.ConnManager
}

//初始化Server模块方法
func NerServer() *Server {
	config := NewConfig()
	s := &Server{
		Name:        config.Name,
		IPversion:   config.Version,
		IP:          config.Host,
		Port:        config.TcpPort,
		MaxConn:     config.MaxConn,
		Routermange: NewRoutermange(),
		ConnManager: NewConnManger(),
		Pool:        NewWorkerPool(2, 10),
	}
	return s
}
