package partition

import (
	"io/ioutil"
	"strconv"

	"github.com/seefan/ssdbproxy/common"
)

func getPid(path string) (processId int, err error) {
	if ex, err := common.FileExists(path); err == nil && ex == true {
		if pid, err := ioutil.ReadFile(path); err == nil {
			processId, err = strconv.Atoi(string(pid))
		}
	}
	return
}