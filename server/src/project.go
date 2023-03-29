package project

import (
	"project/build"
	"project/config"
	"project/zj"

	"github.com/zhengkai/life-go"
)

// Start ...
func Start() {

	build.DumpBuildInfo()

	zj.Init()

	life.Wait()
}

// Prod ...
func Prod() {

	config.Prod = true

	Start()
}
