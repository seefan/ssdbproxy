package server

import (
	"fmt"
	"net"

	log "github.com/cihub/seelog"
	"github.com/seefan/goerr"
	"github.com/seefan/ssdbproxy/conf"
	"github.com/seefan/ssdbproxy/partition"
)

//ssdb代理
type SSDBProxy struct {
	config   *conf.ProxyConf
	listener *net.TCPListener
	isRun    bool
	front    *FrontService
}

func NewSSDBProxy(c *conf.ProxyConf) *SSDBProxy {
	return &SSDBProxy{
		config: c,
		front:  NewFrontService(c.SSDBProxy.PoolSize),
	}
}
func (s *SSDBProxy) Start() error {

	host := fmt.Sprintf("%s:%d", s.config.SSDBProxy.Host, s.config.SSDBProxy.Port)
	tcpAddr, err := net.ResolveTCPAddr("tcp", host)
	if err != nil {
		return goerr.New("resolve tcp address error", host)
	}
	// Listen on TCP port 2000 on all interfaces.
	s.listener, err = net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		return goerr.New("无法启动 SSDBProxy", err)
	}

	log.Infof("SSDBProxy start at %s:%d ", s.config.SSDBProxy.Host, s.config.SSDBProxy.Port)
	s.isRun = true

	go s.run()
	return nil
}
func (s *SSDBProxy) run() {
	for s.isRun {
		conn, err := s.listener.AcceptTCP()
		if err != nil {
			log.Warn("accept conn error", err)
			continue
		}
		s.front.Push(conn)
	}
}
func (s *SSDBProxy) Close() {
	s.isRun = false
	s.listener.Close()
}
