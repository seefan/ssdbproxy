package main

import "github.com/seefan/ssdbproxy/conf"
import (
	"github.com/seefan/ssdbproxy/common"
	"path/filepath"
	"os"
)

func main() {
	defer common.PrintErr()
	baseDir := filepath.Dir(os.Args[0])
	common.InitLog(common.JoinPath(baseDir, "log.xml"), common.JoinPath(baseDir, "logs", "proxy.log"))
	cf := conf.NewProxyConf()
	cf.Load(common.JoinPath(baseDir, "conf.yaml"))
}
