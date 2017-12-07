package partition

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/seefan/goerr"
	"github.com/seefan/ssdbproxy/common"
)

func saveJson(obj interface{}, path string) {
	if bs, err := json.Marshal(obj); err == nil {
		if f, err := os.Create(path); err == nil {
			f.Write(bs)
			f.Close()
		}
	}
}
func readJson(obj interface{}, path string) (err error) {
	if ex, err := common.FileExists(path); ex == true && err == nil {
		if bs, err := ioutil.ReadFile(path); err == nil {
			return json.Unmarshal(bs, obj)
		}
	}
	return goerr.New("port_map.json not found")
}
