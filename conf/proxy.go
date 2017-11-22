package conf

import (
	"gopkg.in/yaml.v2"
	"github.com/seefan/ssdbproxy/common"
	"github.com/seefan/goerr"
	"io/ioutil"
)

type ProxyConf struct {
	SSDBProxy struct {
		Host string `yaml:"host"`
		Port int
		//运行格式
		Model string
		//连接池大小
		PoolSize int `yaml:"pool_size"`
	} `yaml:"ssdb_proxy"`
	SSDB struct {
		Host string
		Port int
		//ssdb的安装目录
		Root string
		//ssdb的工作目录
		Work string
	}
	//分区
	Partition struct {
		//分区模式
		Model string
		//分区限制
		Limit int
		//如何解析key
		Pattern string
	}
}

//新建配置，同时给出默认值
func NewProxyConf() *ProxyConf {
	pc := new(ProxyConf)
	pc.SSDBProxy.Host = "127.0.0.1"
	pc.SSDBProxy.Port = 9999
	pc.SSDBProxy.Model = "cluster"
	pc.SSDBProxy.PoolSize = 10
	pc.SSDB.Host = "127.0.0.1"
	pc.SSDB.Port = 8888
	pc.SSDB.Root = "/usr/local/ssdb"
	pc.SSDB.Work = "/usr/local/ssdb/var"
	pc.Partition.Model = "day"
	pc.Partition.Limit = 5
	pc.Partition.Pattern = "2016-01-02 15"
	return pc
}

//加载配置文件
//
// path string 保存的路径
// 返回 error 加载时的可能错误
func (p *ProxyConf) Load(path string) error {
	exists, err := common.FileNotExist(path)
	if err != nil || !exists {
		return goerr.NewError(err, "ssdb proxy conf error")
	}
	bs, err := ioutil.ReadFile(path)
	if err != nil {
		return goerr.NewError(err, "ssdb proxy conf read error")
	}
	err = yaml.Unmarshal(bs, p)
	if err != nil {
		return goerr.NewError(err, "ssdb proxy conf unmarshal error")
	}
	return nil
}

//保存配置文件
//
// path string 保存的路径
// 返回 error 保存时的可能错误
func (p *ProxyConf) Save(path string) error {
	bs, err := yaml.Marshal(p)
	if err != nil {
		return goerr.NewError(err, "ssdb proxy conf marshal error")
	}
	err = ioutil.WriteFile(path, bs, 0666)
	if err != nil {
		return goerr.NewError(err, "ssdb proxy conf save error")
	}
	return nil
}
