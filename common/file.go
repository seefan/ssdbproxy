package common

import "os"
import "strings"

//文件是否存在，不存在固定返回false，返回true时可能是文件存在，也可能有错误发生。
//
// path string 文件路径
// 返回 bool 文件不存在返回false
// 返回 error 文件存在返回nil，如果有错误返回错误内容
func FileNotExist(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

//合成路径
func JoinPath(path ...string) string {
	return strings.Join(path, string(os.PathSeparator))
}
