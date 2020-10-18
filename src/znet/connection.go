package znet

import (
	"fmt"
	"net"
)

type IConnection interface {
	//启动链接：让当前链接开始工作
	Start()
	//停止链接：结束当前链接的工作
	Stop()
	//获取当前链接的绑定socket conn
	GetTcpConnecttion() *net.TCPConn
	//获取当前链接模块的ID
	GetConnID() uint32
	//获取远程客户端tpc的状态、ip、port
	GetRemoteAddr() net.Addr
	//发送数据、将数据发送给远程的客户端
	Send(data []byte) error
}

//链接模块
type Connection struct {
	//当前链接绑定的服务端
	Server IServer
	//当前链接的socket tcp套接字
	Conn *net.TCPConn
	//链接的ID
	ConnID uint32

	//当前的链接状态
	ConnStatus bool

	//告知当前链接已经退出的/停止的channel
	ExitChan chan bool

	//路由管理器
	Routermange IRoutermanger

	//读协程与写协程之间数据信道
	done chan IRequest
}

func (c *Connection) StartReader() {
	go c.Reader()
	go c.Write()
}

//读
func (c *Connection) Reader() {
	fmt.Println("Reader Goroutine is running")
	defer fmt.Println("connID=", c.ConnID, "Reader is exit,remote addr is", c.GetRemoteAddr().String())
	for {
		buf := make([]byte, 1024)
		_, err := c.Conn.Read(buf)
		if err != nil {
			fmt.Printf("recv buf err %s\n", err)
			c.ExitChan <- true //通知Write，客户端断开链接
			break
		}
		c.ExitChan <- false //如果客户端没有断开链接，向信道发送信息，继续执行解包
		//将数据解包，得到结构体对象
		// msg:封装数据长度、数据id（类型）、数据内容
		msg := NewMessage(buf)
		//将当前链接信息和请求数据封装
		req := NewRequest(c, msg)
		//将解包后的数据传入信道，数据传入信道后会也会阻塞，直到信道中的数据被取走
		c.done <- req
	}
}

//写
func (c *Connection) Write() {
	fmt.Println("Writer Goroutine is running")
	for {
		//信道阻塞，只有从信道中接收到信息才能往下执行
		status := <-c.ExitChan
		if status {
			c.Stop()
			break
		} else {
			req := <-c.done
			//执行注册的路由方法，每次将req传入
			go func(request IRequest) {
				//获取数据类型，查询路由表执行不同的路由
				id := req.Getdata().Getmsgid()
				c.Routermange.IndexRouter(id).Handle(request)
			}(req)
		}
	}
}

//启动链接：让当前链接开始准备工作
func (c *Connection) Start() {
	fmt.Println("Start Conn,ConnID:", c.ConnID)
	//TODO 启动从当前链接写数据的业务
	go c.StartReader()
}

//停止链接：结束当前链接的工作
func (c *Connection) Stop() {
	fmt.Println("Stop Conn,ConnID:", c.ConnID)
	if c.ConnStatus == true {
		return
	}
	//设置状态
	c.ConnStatus = false
	//关闭链接
	c.Conn.Close()
	c.Server.GetConnManager().Remove(c.ConnID)
	//回收资源
	close(c.ExitChan)
	close(c.done)
}

//获取当前链接的绑定socket conn
func (c *Connection) GetTcpConnecttion() *net.TCPConn {
	return c.Conn
}

//获取当前链接模块的ID
func (c *Connection) GetConnID() uint32 {
	return c.ConnID
}

//获取远程客户端tpc的状态、ip、port
func (c *Connection) GetRemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

//发送数据、将数据发送给远程的客户端
func (c *Connection) Send(data []byte) error {
	_, err := c.Conn.Write([]byte(fmt.Sprintf("Hello")))
	return err
}

//此方法用来初始化每一个客户端链接
func NewConnection(Server IServer, conn *net.TCPConn, ConnID uint32, Routermange IRoutermanger) *Connection {
	c := &Connection{
		Server:      Server,
		Conn:        conn,
		ConnID:      ConnID,
		ConnStatus:  false,
		Routermange: Routermange,
		ExitChan:    make(chan bool, 1),
		done:        make(chan IRequest),
	}
	Server.GetConnManager().Add(c.ConnID, c)
	return c
}
