package zj

import (
	"path/filepath"
	"project/config"

	"github.com/zhengkai/zog"
)

// Init ...
func Init() {

	mainCfg := zog.NewConfig()
	mainCfg.Caller = zog.CallerLong

	infoCfg := mainCfg.Clone()
	infoCfg.Color = zog.ColorInfo
	infoCfg.LinePrefix = `[IO] `

	debugCfg := mainCfg.Clone()
	debugCfg.Color = zog.ColorLight
	debugCfg.LinePrefix = `[Debug] `

	errCfg := zog.NewErrConfig()
	errCfg.Color = zog.ColorWarn
	errCfg.LinePrefix = `[Error] `

	baseLog.CDefault = mainCfg
	baseLog.CDebug = debugCfg
	baseLog.CInfo = infoCfg
	baseLog.CError = errCfg
	baseLog.CWarn = errCfg
	baseLog.CFatal = errCfg

	baseLog.SetDirPrefix(filepath.Dir(zog.GetSourceFileDir()))

	dir := config.LogDir

	mainFile, _ := zog.NewFile(dir+`/default.txt`, false)
	infoFile, _ := zog.NewFile(dir+`/io.txt`, false)
	errFile, _ := zog.NewFile(dir+`/err.txt`, true)

	mainCfg.Output = append(mainCfg.Output, mainFile)
	infoCfg.Output = append(infoCfg.Output, infoFile)
	errCfg.Output = append(errCfg.Output, mainFile, errFile)
}
