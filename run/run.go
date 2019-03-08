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

func Run(r server.Runnable, outFile ...string) {
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
	pidFile := filepath.Join(path, "pid.save")

	switch cmd {
	case "start":
		if err := savePid(pidFile); err != nil {
			println(err.Error())
		}
		var f *os.File
		if len(outFile) > 0 {
			logFile := filepath.Join(path, outFile[0])
			tmp := filepath.Dir(logFile)
			if _, err := os.Stat(tmp); os.IsNotExist(err) {
				if err := os.MkdirAll(tmp, 0764); err != nil {
					panic(err.Error())
				}
			}
			if f, err = os.Create(logFile); err != nil {
				panic(err.Error())
			}
			defer func() {
				if err = f.Close(); err != nil {
					println(err.Error())
				}
			}()
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
		}, f, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

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
