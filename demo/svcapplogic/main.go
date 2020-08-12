package main

import (
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"iogo/demo/comm/iogo"
	_ "iogo/demo/svcapplogic/server"
)

const defServicePort = "28667"

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU() * 2)

	isError := true
	waitStop := func() {
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL, syscall.SIGHUP, syscall.SIGQUIT)
		go func() {
			<-ch
			isError = false
			iogo.StopIogo()
			os.Exit(1)
		}()
	}
	waitStop()

	iogo.RunIogo(defServicePort)

	if isError {
		iogo.StopIogo()
	}
}
