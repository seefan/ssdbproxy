package partition

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"

	"github.com/seefan/goerr"
	"github.com/seefan/ssdbproxy/common"
)

func saveSSDBConfig(host, work, root string, port int) error {
	ssdbConf := common.JoinPath(root, "ssdb.conf")
	outConf := common.JoinPath(work, "ssdb.conf")

	ne, err := common.FileExists(ssdbConf)
	if ne == false || err != nil {
		return goerr.New("ssdb.conf not found", ssdbConf)
	}
	f, err := os.Open(ssdbConf)
	if err != nil {
		return err
	}
	defer f.Close()

	fo, err := os.Create(outConf)
	if err != nil {
		return err
	}
	defer fo.Close()
	out := bufio.NewWriter(fo)
	cfg := bufio.NewReader(f)

	//find work_dir
	if b, e := find(cfg, []byte("work_dir")); e == nil {
		out.Write(b)
		out.WriteString("work_dir=" + work + "\n")
	} else {
		return e
	}
	//find work_dir
	if b, e := find(cfg, []byte("pidfile")); e == nil {
		out.Write(b)
		out.WriteString("pidfile=" + common.JoinPath(work, "ssdb.pid") + "\n")
	} else {
		return e
	}
	//find ip
	if b, e := find(cfg, []byte("ip")); e == nil {
		out.Write(b)
		out.WriteString("\tip:" + host + "\n")
	} else {
		return e
	}
	//find port
	if b, e := find(cfg, []byte("port")); e == nil {
		out.Write(b)
		out.WriteString(fmt.Sprintf("\tport:%d\n", port))
	} else {
		return e
	}
	//find log output
	if b, e := find(cfg, []byte("output")); e == nil {
		out.Write(b)
		out.WriteString("\toutput:" + common.JoinPath(work, "logs", "log.txt") + "\n")
	} else {
		return e
	}
	for {
		bs, _, e := cfg.ReadLine()
		if e == io.EOF {
			break
		}
		if e != nil {
			return e
		}
		out.Write(bs)
		out.WriteRune('\n')
	}
	out.Flush()
	return nil
}
func find(reader *bufio.Reader, str []byte) ([]byte, error) {
	var out []byte
	for {
		bs, _, e := reader.ReadLine()
		if e == io.EOF {
			break
		}
		if e != nil {
			return nil, e
		}
		bb := bytes.TrimSpace(bs)
		if !bytes.HasPrefix(bb, []byte{'#'}) {
			if bytes.Index(bs, str) != -1 {
				return out, nil
			}
		}

		out = append(out, bs...)
		out = append(out, '\n')
	}
	return nil, goerr.New("not found")
}
