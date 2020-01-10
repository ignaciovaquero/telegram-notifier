package controllers

import "log"

// Logger defines an interface for logs
type Logger interface {
	// logging facilities
	Debug(...interface{})
	Debugf(string, ...interface{})
	Error(...interface{})
	Errorf(string, ...interface{})
	Info(...interface{})
	Infof(string, ...interface{})
}

type defaultLogger struct {
	*log.Logger
}

func (d *defaultLogger) Debug(v ...interface{}) {
	d.Logger.Print(v...)
}

func (d *defaultLogger) Debugf(format string, v ...interface{}) {
	d.Logger.Printf(format, v...)
}

func (d *defaultLogger) Error(v ...interface{}) {
	d.Logger.Print(v...)
}

func (d *defaultLogger) Errorf(format string, v ...interface{}) {
	d.Logger.Printf(format, v...)
}

func (d *defaultLogger) Info(v ...interface{}) {
	d.Logger.Print(v...)
}

func (d *defaultLogger) Infof(format string, v ...interface{}) {
	d.Logger.Printf(format, v...)
}
