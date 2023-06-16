package project

import (
	"project/build"
	"project/config"
	"project/es"
	"project/web"
	"project/zj"

	"github.com/zhengkai/life-go"
)

// Start ...
func Start() {

	build.DumpBuildInfo()

	zj.Init()

	es.Init()

	if !config.Prod {
		// es.LastItem()
		// st.DateHistogram()
		go es.Test()
	}

	go web.Server()

	life.Wait()
}

// Prod ...
func Prod() {

	config.Prod = true

	Start()
}
