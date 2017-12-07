package partition

import (
	"io/ioutil"
	"os"
	"os/exec"
	"time"

	log "github.com/cihub/seelog"
	"github.com/seefan/ssdbproxy/common"
	"github.com/seefan/ssdbproxy/conf"
)

var (
	watcher *time.Ticker
	portMap = make(map[int]interface{})
)

func Close() {
	if watcher != nil {
		watcher.Stop()
	}
}
func Start(config *conf.ProxyConf) error {
	if err := loadSSDB(config, time.Now()); err != nil {
		return err
	}
	watcher = time.NewTicker(time.Minute)
	go watchSSDB(config, watcher.C)
	return nil
}
func watchSSDB(config *conf.ProxyConf, dt <-chan time.Time) {
	for t := range dt {
		if t.Minute()%60 == 0 { //10分钟执行一次清理检查
			//获取所有ssdb目录
			if dirList, e := ioutil.ReadDir(config.SSDB.Work); e == nil {
				dirs := pt.Clean(config.Partition.Limit, dirList)
				for _, dir := range dirs {
					log.Infof("clean dir %s", dir)
					//去掉port映射
					for port, d := range portMap {
						if d == dir {
							portMap[port] = nil
						}
					}
					//关ssdb
					cmd := exec.Command(common.JoinPath(config.SSDB.Root, "ssdb-server"), common.JoinPath(config.SSDB.Work, dir, "ssdb.conf"), "-s", "stop")
					_, err := cmd.CombinedOutput()
					if err == nil {
						//删除目录
						if err := os.RemoveAll(common.JoinPath(config.SSDB.Work, dir)); err != nil {
							log.Warnf("remove dir(%s) error", dir, err)
						}
					}
				}
			}
			portMapFile := common.JoinPath(config.SSDB.Work, "port_map.json")
			saveJson(portMap, portMapFile)
		}
		loadSSDB(config, t) //一分钟检查一个ssdb的状态
	}
}
func findPort(dir string) int {
	for k, v := range portMap {
		if v == nil || v == dir {
			return k
		}
	}
	return -1
}
func loadSSDB(config *conf.ProxyConf, start time.Time) error {
	tomorrow := start.AddDate(0, 0, 1)
	//初始化port和ssdb的对应关系
	portMapFile := common.JoinPath(config.SSDB.Work, "port_map.json")
	portMapTmp := make(map[int]interface{})
	if err := readJson(&portMapTmp, portMapFile); err == nil {
		for k, v := range portMapTmp {
			portMap[k] = v
		}
	} else {
		//init pool map
		for i := 0; i < config.Partition.Limit; i++ {
			portMap[config.SSDB.Port+i] = nil
		}
	}
	var allKey []string
	switch config.Partition.Model {
	case "day":
		t := NewDayPartition(config)
		allKey = t.getAllKey(tomorrow, config.Partition.Limit)
		pt = t
	}
	for i, ssdbKey := range allKey {
		log.Infof("load ssdb %d,key is %s", i, ssdbKey)
		//mainKey = now.Format(config.Partition.Pattern)
		work := common.JoinPath(config.SSDB.Work, ssdbKey)
		ex, err := common.FileExists(work)
		if err != nil {
			return err
		}
		if ex == false {
			os.MkdirAll(work, 0764)
		}
		logFile := common.JoinPath(work, "logs")
		ex, err = common.FileExists(logFile)
		if err != nil {
			return err
		}
		if ex == false {
			os.MkdirAll(logFile, 0764)
		}

		port := findPort(ssdbKey)
		if port == -1 {
			return nil
		}
		//set ssdb config
		if err := saveSSDBConfig(config.SSDB.Host, work, config.SSDB.Root, port); err != nil {
			return err
		}
		time.Sleep(time.Second)
		//start ssdb
		if err := startSSDB(config.SSDB.Root, work); err != nil {
			return err
		}
		//create ssdb pool
		//pool, err := gossdb.NewPool(&sp.Config{
		//	Host: config.SSDB.Host,
		//	Port: port,
		//})
		//if err != nil {
		//	return err
		//}
		//poolMap[ssdbKey] = pool
		portMap[port] = ssdbKey
	}
	saveJson(portMap, portMapFile)
	return nil
}
