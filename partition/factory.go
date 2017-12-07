package partition

import (
	"os"

	"github.com/seefan/goerr"
	"github.com/seefan/gossdb"
)

var (
	pt      IPartition
	poolMap = make(map[string]*gossdb.Connectors)
	//keyPort = make(map[string]int)
	mainKey string
)

func Get(key string) (*gossdb.Client, error) {
	if pt != nil {
		if key == "" {
			key = mainKey
		}
		if k, err := pt.Get(key); err == nil {
			if c, ok := poolMap[k]; ok {
				return c.NewClient()
			}
		}
	}
	return nil, goerr.New("key is wrong", key)
}

type IPartition interface {
	Get(string) (string, error)
	Clean(int, []os.FileInfo) []string
}
