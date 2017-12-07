package common

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	log "github.com/cihub/seelog"
)

var (
	logConfig = `<seelog type="asynctimer" asyncinterval="5000000" minlevel="debug">
	<outputs formatid="main">
		<console/>
		<rollingfile type="size" filename="__log_url__" maxsize="1024000" maxrolls="10" />
	</outputs>
	<formats>
		<format id="main" format="%Date(2006-01-02 15:04:05) [%Level] %RelFile line:%Line %Msg%n"/>
	</formats>
</seelog>`
)

func InitLog(configFile, logFile string) {
	if exists, _ := FileExists(configFile); exists == false {
		ioutil.WriteFile(configFile, []byte(strings.Replace(logConfig, "__log_url__", logFile, 1)), 0666)
	}
	if logger, err := log.LoggerFromConfigAsFile(configFile); err == nil {
		log.ReplaceLogger(logger)
	}
}
func PrintErr() {
	if err := recover(); err != nil {
		path, fe := filepath.Abs(os.Args[0])
		if fe != nil {
			path = os.Args[0]
		}
		path = filepath.Dir(path)
		path += string(os.PathSeparator) + "fault.txt"
		str := fmt.Sprintf("%v\n", err)
		for i := 0; i < 10; i++ {
			funcName, file, line, ok := runtime.Caller(i)
			if ok {
				str += fmt.Sprintf("frame %v:[func:%v,file:%v,line:%v]\n", i, runtime.FuncForPC(funcName).Name(), file, line)
			}
		}
		println(str)
		ioutil.WriteFile(path, []byte(str), 0764)
	}
}
