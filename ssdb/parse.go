package ssdb

import (
	"bytes"
	"net"
	"strconv"

	"github.com/seefan/goerr"
)

const (
	ENDN = '\n'
	ENDR = '\r'
)

var (
	byt              = []int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	maxByteSize byte = 57
	minByteSize byte = 48
)

func ToNum(bs []byte) int {
	re := 0
	for _, v := range bs {
		if v > maxByteSize || v < minByteSize {
			return re
		}
		re = re*10 + byt[v]
	}
	return re
}
func Encode(c net.Conn, args []string) error {
	var packetBuf bytes.Buffer
	for _, arg := range args {
		packetBuf.Write(strconv.AppendInt(nil, int64(len(arg)), 10))
		packetBuf.WriteByte(ENDN)
		packetBuf.WriteString(arg)
		packetBuf.WriteByte(ENDN)
	}
	packetBuf.WriteByte(ENDN)

	for _, err := packetBuf.WriteTo(c); packetBuf.Len() > 0; {
		if err != nil {
			packetBuf.Reset()
			return goerr.NewError(err, "client socket write error")
		}
	}

	packetBuf.Reset()
	return nil
}
func Decode(c net.Conn) (resp []interface{}, err error) {

	bufSize := 0
	var packetBuf []byte
	readBuf := make([]byte, 8*1024)
	//数据包分解，发现长度，找到结尾，循环发现，发现空行，结束
	for {
		bufSize, err = c.Read(readBuf)
		if err != nil {
			return nil, goerr.NewError(err, "client socket read error")
		}
		if bufSize < 1 {
			continue
		}
		packetBuf = append(packetBuf, readBuf[:bufSize]...)

		for {
			rsp, n := parse(packetBuf)
			if n == -1 {
				break
			} else if n == -2 {
				return
			} else {
				resp = append(resp, rsp)
				packetBuf = packetBuf[n+1:]
			}
		}
	}
	packetBuf = nil

	return resp, nil
}

//解析数据为string的slice
func parse(buf []byte) (resp string, size int) {
	n := bytes.IndexByte(buf, ENDN)
	blockSize := -1
	size = -1
	if n != -1 {
		if n == 0 || n == 1 && buf[0] == ENDR { //空行，说明一个数据包结束
			size = -2
			return
		}
		//数据包开始，包长度解析
		blockSize = ToNum(buf[:n])
		bufSize := len(buf)

		if n+blockSize < bufSize {
			resp = string(buf[n+1: blockSize+n+1])
			for i := blockSize + n + 1; i < bufSize; i++ {
				if buf[i] == ENDN {
					size = i
					return
				}
			}
		}
	}

	return
}
