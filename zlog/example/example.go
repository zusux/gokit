package main

import (
	"github.com/zusux/gokit/zlog"
)

func main() {
	(&zlog.Logger{}).InitLog()
	zlog.Info("this is info")
	zlog.Warn("this is warn")
	zlog.Error("this is error1")
	zlog.Error("this is error2")
	zlog.Infof("%s xxx", "fengfeng")
}
