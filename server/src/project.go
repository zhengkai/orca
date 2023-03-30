package project

import (
	"project/build"
	"project/config"
	"project/web"
	"project/zj"

	"github.com/zhengkai/life-go"
)

// Start ...
func Start() {

	build.DumpBuildInfo()

	zj.Init()

	go web.Server()

	life.Wait()
}

// Prod ...
func Prod() {

	config.Prod = true

	Start()
}
