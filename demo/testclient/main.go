package main

import (
	"bufio"
	"fmt"
	"os"
	"runtime"

	"iogo/demo/comm/iogo"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU() * 2)

	iogo.RunIogo("")

	testStart()

	fmt.Println("Press the enter key to exit")
	reader := bufio.NewReader(os.Stdin)
	reader.ReadString('\n')

	iogo.StopIogo()
}
