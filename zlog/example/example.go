package main

import (
	"github.com/zusux/gokit/zlog"
	"go.uber.org/zap"
)

func main() {
	logger := zlog.NewLogger().InitLog()
	logger.Info("this is info")
	logger.Warn("this is warn")
	logger.Error("this is error1")
	logger.Error("this is error2")
	zap.S().Infof("%s xxx", "fengfeng")
}
