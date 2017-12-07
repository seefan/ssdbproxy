package partition

import (
	"os"
	"os/exec"

	log "github.com/cihub/seelog"
	"github.com/seefan/ssdbproxy/common"
)

func startSSDB(root, work string) error {
	if checkProcess(work) {
		return nil
	}
	log.Infof("start ssdb @ %s",work)
	cmd := exec.Command(common.JoinPath(root, "ssdb-server"), common.JoinPath(work, "ssdb.conf"), "-d", "-s", "restart")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}
	
	log.Trace(string(out))
	return nil
}
func checkProcess(work string) bool {
	pidFile := common.JoinPath(work, "ssdb.pid")
	if processId, err := getPid(pidFile); err == nil {
		if p, err := os.FindProcess(processId); err == nil {
			return p.Pid > 0
		} else {
			log.Error(err)
		}
	}

	return false
}
func cleanSSDB(root, work string) error {
	cmd := exec.Command(common.JoinPath(root, "ssdb-server"), common.JoinPath(work, "ssdb.conf"), "-s", "stop")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}
	log.Trace(string(out))
	if checkProcess(work) {
		cmd := exec.Command(common.JoinPath(root, "ssdb-server"), common.JoinPath(work, "ssdb.conf"), "-s", "stop")
		_, err := cmd.CombinedOutput()
		if err != nil {
			return err
		}
	}
	if checkProcess(work) == false {
		return os.RemoveAll(work)
	}
	return nil
}
