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

	Task("build-mobile").
		Chdir("cmd/mobifire").
		Exec("fyne", "package", "-os", "android", "-appID", "net.kettek.mobifire", "-icon", "icon.png")

	Task("run-mobile").
		Exec("adb", "shell", "am", "force-stop", "net.kettek.mobifire").
		Exec("adb", "shell", "am", "start", "-a", "android.intent.action.MAIN", "-n", "net.kettek.mobifire/org.golang.app.GoNativeActivity")

	Task("install-mobile").
		Chdir("cmd/mobifire").
		Exec("adb", "install", "./mobifire.apk")

	Task("watch-mobile").
		Watch("states/**/*.go", "net/**/*.go").
		Run("build-mobile").
		Run("install-mobile").
		Run("run-mobile").
		Sleep("1y")

	Go()
}
