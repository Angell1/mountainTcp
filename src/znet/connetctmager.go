package znet

import (
	"errors"
	"sync"
)

type IConnManager interface {
	Add(connID uint32, conn IConnection)
	Remove(connID uint32)
	GetConn(connID uint32) (IConnection, error)
	GetConnLen() int
	ClearConn()
}
type ConnManger struct {
	ConnMangerMap map[uint32]IConnection
	connLock      *sync.RWMutex
}

//添加
func (M *ConnManger) Add(connID uint32, conn IConnection) {
	M.connLock.Lock()         //加写锁
	defer M.connLock.Unlock() //解写锁
	M.ConnMangerMap[connID] = conn

}

//删除
func (M *ConnManger) Remove(connID uint32) {
	M.connLock.Lock()         //加写锁
	defer M.connLock.Unlock() //解写锁
	delete(M.ConnMangerMap, connID)
}

//得到链接
func (M *ConnManger) GetConn(connID uint32) (IConnection, error) {
	M.connLock.RLock()         //加读锁
	defer M.connLock.RUnlock() //解读锁
	if conn, ok := M.ConnMangerMap[connID]; ok {
		return conn, nil
	} else {
		return nil, errors.New("connnection not Found!")
	}

}

//得到链接个数
func (M *ConnManger) GetConnLen() int {
	return len(M.ConnMangerMap)
}

func (M *ConnManger) ClearConn() {
	M.connLock.Lock()         //加写锁
	defer M.connLock.Unlock() //解写锁
	for ConnID, conn := range M.ConnMangerMap {
		//停止链接
		conn.Stop()
		//删除链接
		delete(M.ConnMangerMap, ConnID)
	}
}

func NewConnManger() *ConnManger {
	return &ConnManger{
		make(map[uint32]IConnection),
		new(sync.RWMutex),
	}
}
