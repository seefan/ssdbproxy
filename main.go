package main

import (
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	log "github.com/cihub/seelog"
	"github.com/seefan/ssdbproxy/common"
	"github.com/seefan/ssdbproxy/conf"
	"github.com/seefan/ssdbproxy/partition"
)

func main() {
	defer common.PrintErr()
	baseDir := filepath.Dir(os.Args[0])
	common.InitLog(common.JoinPath(baseDir, "log.xml"), common.JoinPath(baseDir, "logs", "proxy.log"))
	defer log.Flush()
	cf := conf.NewProxyConf()
	cf.Load(common.JoinPath(baseDir, "conf.yaml"))
	log.Info("SSDBProxy starting")
	//proxy := server.NewSSDBProxy(cf)
	//if err := proxy.Start(); err != nil {
	//	log.Critical(err)
	//}
	if err := partition.Start(cf); err != nil {
		log.Error("partition init error", err)
	}
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, os.Kill, syscall.SIGTERM, syscall.SIGINT)
	<-interrupt
	//proxy.Close()
	log.Info("SSDBProxy closed")
}
