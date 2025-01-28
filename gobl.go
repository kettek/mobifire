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

	runArgs := append([]interface{}{}, "./mobifire"+exe)

	Task("build").
		Exec("go", "build", "./cmd/mobifire")

	Task("run").
		Exec(runArgs...)

	Task("watch").
		Watch("**/*.go").
		Signaler(SigQuit).
		Run("build").
		Run("run")

	Go()
}
