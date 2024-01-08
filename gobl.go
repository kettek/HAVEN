package main

import (
	"runtime"

	. "github.com/kettek/gobl"
)

func main() {
	var exe string
	if runtime.GOOS == "windows" {
		exe = ".exe"
	}

	runArgs := append([]interface{}{}, "./ebihack23"+exe)

	Task("build").
		Exec("go", "build", "./cmd/ebihack23")
	Task("run").
		Exec(runArgs...)
	Task("watch").
		Watch("./cmd/*/*.go", "./res/*", "./game/*.go", "./rooms/*.go", "./states/*.go", "./actors/*.go", "./inputs/*.go", "./commands/*.go").
		Signaler(SigQuit).
		Run("build").
		Run("run")
	Go()
}
