package znet

type IRequest interface {
	//得到请求数据封装
	Getdata() IMessage
	//得到当前链接封装
	GetConnecttion() IConnection
}

type Request struct {
	//已经和客户端建立链接对象
	conn IConnection
	//客户端请求数据
	Msg IMessage
}

//得到请求数据
func (r *Request) Getdata() IMessage {
	return r.Msg
}

//得到当前链接
func (r *Request) GetConnecttion() IConnection {
	return r.conn
}

//初始化Request对象
func NewRequest(conn IConnection, Msg IMessage) *Request {
	r := &Request{
		conn: conn,
		Msg:  Msg,
	}
	return r
}
