/*
@Time : 2019-01-25 15:55
@Author : seefan
@File : run
@Software: microgo
*/
package run

import (
	"context"
	"github.com/seefan/microgo/server"
	"os"
	"path/filepath"
	"syscall"
)

func Run(r server.Runnable, outFile ...*os.File) {
	defer printErr()
	cmd := "debug"
	if len(os.Args) > 1 {
		cmd = os.Args[1]
	}

	path, err := filepath.Abs(os.Args[0])
	if err != nil {
		path = os.Args[0]
	}
	path = filepath.Dir(path)
	path += string(os.PathSeparator)

	pidFile := path + "pid.save"

	switch cmd {
	case "start":
		if err := savePid(pidFile); err != nil {
			println(err.Error())
		}
		var f *os.File
		if len(outFile) > 0 {
			f = outFile[0]
		}
		nohup(func() error {
			return start(path, r)
		}, func(signal os.Signal, e error) {
			if e != nil {
				println(e.Error())
			} else {
				if e = r.Stop(); e != nil {
					println(e.Error())
				}
			}
		}, f, syscall.SIGINT, syscall.SIGKILL)
	case "stop":
		if err := stopByPidFile(pidFile); err != nil {
			println(err.Error())
		}
	default:
		if err := start(path, r); err != nil {
			println(err.Error())
		}
	}
}

func start(startPath string, r server.Runnable) error {
	//load config

	if err := r.Start(context.Background()); err != nil {
		return err
	}
	return nil
}
