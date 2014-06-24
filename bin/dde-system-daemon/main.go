package main

import "pkg.linuxdeepin.com/lib/logger"

import "pkg.linuxdeepin.com/lib"
import "pkg.linuxdeepin.com/lib/dbus"
import "os"
import _ "dde-daemon/accounts"
import "dde-daemon"

var Logger = logger.NewLogger("com.deepin.daemon")

func main() {
	Logger.BeginTracing()
	defer Logger.EndTracing()

	if !lib.UniqueOnSystem("com.deepin.daemon") {
		Logger.Warning("There already has an dde-daemon running.")
		return
	}

	Logger.SetRestartCommand("/usr/lib/deepin-daemon/dde-system-daemon")

	loader.Start()
	defer loader.Stop()

	dbus.DealWithUnhandledMessage()

	if err := dbus.Wait(); err != nil {
		Logger.Errorf("Lost dbus: %v", err)
		os.Exit(-1)
	} else {
		os.Exit(0)
	}
}
