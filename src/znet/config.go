package znet

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"runtime"
	"strconv"
)

type Config struct {
	Name        string
	Version     string
	Host        string
	TcpPort     string
	MaxPackSize int64
	MaxConn     int64
}

func (c *Config) setname(Name string) {
	c.Name = Name
}
func (c *Config) setVersion(Version string) {
	c.Version = Version
}

func (c *Config) setHost(Host string) {
	c.Host = Host
}

func (c *Config) setTcpPort(TcpPort string) {
	c.TcpPort = TcpPort
}
func (c *Config) setMaxConn(MaxConn int64) {
	c.MaxConn = MaxConn
}

func (c *Config) setMaxPackSize(MaxPackSize int64) {
	c.MaxPackSize = MaxPackSize
}

func NewConfig() *Config {
	//打开json文件
	const dataFile = "./server.json"
	_, filename, _, _ := runtime.Caller(1)
	datapath := path.Join(path.Dir(filename), dataFile)
	srcFile, err := os.Open(datapath)
	if err != nil {
		fmt.Println("文件打开失败,err=", err)
		return nil
	}
	defer srcFile.Close()
	//创建接收数据的map类型数据
	datamap := make(map[string]string)
	//创建解码器
	decoder := json.NewDecoder(srcFile)
	//解码
	err = decoder.Decode(&datamap)
	if err != nil {
		fmt.Println("解码失败,err:", err)
		return nil
	}
	MaxConn, err := strconv.ParseInt(datamap["MaxConn"], 10, 64)
	MaxPackSize, err := strconv.ParseInt(datamap["MaxPackSize"], 10, 64)
	c := &Config{
		datamap["Name"],
		datamap["Version"],
		datamap["Host"],
		datamap["TcpPort"],
		MaxPackSize,
		MaxConn,
	}
	return c
}
