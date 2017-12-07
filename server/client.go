package server

import (
	"net"

	"github.com/seefan/gopool"
	"github.com/seefan/ssdbproxy/partition"

	log "github.com/cihub/seelog"
	"github.com/seefan/ssdbproxy/ssdb"
)

type FrontClient struct {
	conn     chan *net.TCPConn
	callback func(*FrontClient)
	client   *gopool.PooledClient
	isOpen   bool
}

func NewPartClient(f func(client *FrontClient)) *FrontClient {
	return &FrontClient{
		conn:     make(chan *net.TCPConn, 1),
		callback: f,
	}
}

func (p *FrontClient) Push(c *net.TCPConn) {
	p.conn <- c
}

//打开连接
//
// 返回，error。如果连接到服务器时出错，就返回错误信息，否则返回nil
func (p *FrontClient) Start() error {
	p.isOpen = true
	go p.run()
	return nil
}
func (p *FrontClient) run() {
	var key string
	for c := range p.conn {
		for {
			resp, err := ssdb.Decode(c)
			if err != nil {
				break
			}
			key = ssdb.GetKey(resp)
			if sc, e := partition.Get(key); e == nil {
				if resp, err := sc.Do(resp...); err == nil {
					ssdb.Encode(c, resp)
				} else {
					log.Error(err)
				}
			} else {
				log.Error(e)
			}
		}
		p.callback(p)
	}
}

//关闭连接
//
// 返回，error。如果关闭连接时出错，就返回错误信息，否则返回nil
func (p *FrontClient) Close() error {
	p.isOpen = false
	close(p.conn)
	return nil
}

//是否打开
//
// 返回，bool。如果已连接到服务器，就返回true。
func (p *FrontClient) IsOpen() bool {
	return p.isOpen
}

//检查连接状态
//
// 返回，bool。如果无法访问服务器，就返回false。
func (p *FrontClient) Ping() bool {
	return true
}
