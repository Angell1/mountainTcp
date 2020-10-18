package znet

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

type IDatapack interface {
	//获取包的头长度方法
	GetHead() uint32
	//封包
	Pack(msg IMessage) ([]byte, error)
	//解包
	Unpack(data []byte) (Message, error)
}

type Datapack struct{}

func (dp *Datapack) GetHead() uint32 {
	//4字节长度+4字节长度
	return 8

}

func (dp *Datapack) Pack(msg IMessage) ([]byte, error) {
	//创建一个存放bytes的缓冲
	dataBuff := bytes.NewBuffer([]byte{})
	//将数据长度写入字节流
	err := binary.Write(dataBuff, binary.LittleEndian, msg.Getmsglen())
	checkerr(err)
	//将id写入字节流
	err = binary.Write(dataBuff, binary.LittleEndian, msg.Getmsgid())
	checkerr(err)
	//将数据内容写入字节流
	err = binary.Write(dataBuff, binary.LittleEndian, msg.Getmsg())
	checkerr(err)
	return dataBuff.Bytes(), nil

}
func (dp *Datapack) Unpack(data []byte) (*Message, error) {
	//这里可以不需要额外创建一个数据缓冲
	//创建一个io。Reader
	boolBuffer := bytes.NewReader(data)
	msg := &Message{}
	//读取数据长度和id
	err := binary.Read(boolBuffer, binary.LittleEndian, &msg.Msglen)
	checkerr(err)
	err = binary.Read(boolBuffer, binary.LittleEndian, &msg.Msgid)
	checkerr(err)
	//数据包限制
	//if
	//
	//}
	return msg, nil
}

func NewDatapack() *Datapack {
	return &Datapack{}
}

func checkerr(err error) {
	if err != nil {
		fmt.Println("数据写入与读取失败")
	}
}
