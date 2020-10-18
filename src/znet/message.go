package znet

import (
	"fmt"
)

type IMessage interface {
	//获取消息的内容
	Getmsg() []byte
	//获取消息的长度
	Getmsglen() uint32
	//获取消息的ID
	Getmsgid() uint32
	//设置消息的ID
	Setmsgid(id uint32)
	//设置消息的内容
	Setmsg(data []byte)
}

type Message struct {
	//TLV格式数据
	//消息长度
	Msglen uint32
	//消息ID
	Msgid uint32
	//消息内容
	data []byte
}

func (M *Message) Getmsg() []byte {
	return M.data
}
func (M *Message) Getmsglen() uint32 {
	return M.Msglen
}
func (M *Message) Getmsgid() uint32 {
	return M.Msgid
}
func (M *Message) Setmsgid(id uint32) {
	M.Msgid = id
}
func (M *Message) Setmsg(data []byte) {
	M.data = data
}

func NewMessage(data []byte) *Message {
	pack := NewDatapack()
	msg, _ := pack.Unpack(data)
	fmt.Printf("recv buf:%s\n", string(data[8:8+msg.Msglen]))
	msg.data = data[8 : 8+msg.Msglen]
	return msg
}
