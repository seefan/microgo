package run

import (
	"os"
)

func stopByPidFile(pidFile string) error {
	if pid, err := getPid(pidFile); err == nil {
		if err = kill(pid); err != nil {
			return err
		} else {
			return os.Remove(pidFile)
		}
	} else {
		return err
	}
}
