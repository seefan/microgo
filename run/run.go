/*
@Time : 2019-01-25 15:55
@Author : seefan
@File : run
@Software: microgo
*/
package run

import (
	"context"
	"os"
	"path/filepath"
	"syscall"

	"github.com/seefan/microgo/server"
)

//Run the new server.Runnable
func Run(run func() server.Runnable, outFile ...string) {
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
		nohup(func(sig chan<- os.Signal) {
			r := run()
			if err := r.Start(context.Background()); err != nil {
				sig <- syscall.SIGABRT
				return
			}
			if err := r.Stop(); err != nil {
				println(err.Error())
			}

		}, f, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	case "stop":
		if err := stopByPidFile(pidFile); err != nil {
			println(err.Error())
		}
	default:
		r := run()
		if err := r.Start(context.Background()); err != nil {
			println(err.Error())
		}
	}
}
