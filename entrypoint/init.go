package entrypoint

import (
  "fmt"
	. "user/version"
	. "user/global"

	"github.com/spf13/viper"
)

func Init() {
	EntryPoint()
	updateversion()
}

func updateversion() {
setingskey := fmt.Sprintf("%s.entrypointversion", MicroServiceName)
viper.Set(setingskey, VERSIONPLUGIN)
}
