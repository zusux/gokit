package zlog

type SLog struct{}

func NewSLog() *SLog {
	return &SLog{}
}
func (s SLog) Debugf(format string, args ...interface{}) {
	Debugf(format, args...)
}

func (s SLog) Printf(format string, args ...interface{}) {
	Printf(format, args...)
}

func (s SLog) Infof(format string, args ...interface{}) {
	Infof(format, args...)
}
func (s SLog) Warnf(format string, args ...interface{}) {
	Warnf(format, args...)
}
func (s SLog) Errorf(format string, args ...interface{}) {
	Errorf(format, args...)
}
func (s SLog) Panicf(format string, args ...interface{}) {
	Panicf(format, args...)
}
func (s SLog) Fatalf(format string, args ...interface{}) {
	Fatalf(format, args...)
}
