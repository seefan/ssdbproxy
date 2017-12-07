package server

import (
	"net"

	"github.com/gpmgo/gopm/modules/log"
	"github.com/seefan/gopool"
)

//处理前端业务，主要是解析ssdb请求，并发送给ssdb
type FrontService struct {
	pool *gopool.Pool
}

func NewFrontService(count int) *FrontService {
	p := gopool.NewPool()
	p.MaxPoolSize = count

	f := &FrontService{
		pool: p,
	}
	p.NewClient = func() gopool.IClient {
		return NewPartClient(f.callback)
	}
	p.Start()
	return f
}
func (f *FrontService) Push(c *net.TCPConn) {
	var thc *FrontClient
	if th, err := f.pool.Get(); err == nil {
		thc = th.Client.(*FrontClient)
		thc.client = th
		thc.Push(c)
	} else {
		log.Error("connection is full")
		c.Close()
	}
}
func (f *FrontService) callback(c *FrontClient) {
	f.pool.Set(c.client)
}
