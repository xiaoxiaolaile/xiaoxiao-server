package core

import logs "github.com/sirupsen/logrus"

type Console struct {
}

func (c *Console) Info(args ...interface{}) {
	logs.Info(args...)
}
func (c *Console) Debug(args ...interface{}) {
	logs.Debug(args...)
}
func (c *Console) Warn(args ...interface{}) {
	logs.Warn(args...)
}
func (c *Console) Error(args ...interface{}) {
	logs.Error(args...)
}
func (c *Console) Log(format string, args ...interface{}) {
	c.LogF(format, args...)
}
func (c *Console) LogF(format string, args ...interface{}) {
	logs.Printf(format, args...)
}
