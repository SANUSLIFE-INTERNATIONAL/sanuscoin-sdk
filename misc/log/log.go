package log

import (
	"fmt"
	"io/ioutil"
	"os"

	"sanus/sanus-sdk/config"

	"github.com/sirupsen/logrus"
)

type Logger struct {
	logger *logrus.Logger

	prefix string

	out *os.File
}

func NewLogger(cfg *config.Config) *Logger {
	logger := logrus.New()
	logger.Formatter = &logrus.TextFormatter{
		ForceColors:     true,
		DisableColors:   false,
		ForceQuote:      true,
		DisableQuote:    true,
		FullTimestamp:   true,
		TimestampFormat: "Jan_2 / 15:04:05",
	}
	logger.Level = parseLogLevel(cfg.App.Debug)
	l := &Logger{logger: logger}
	return l
}

func (logger *Logger) SetOutput(fName, prefix string) {
	logger.setOut(config.AppLogPath(fName))
	logger.setHook()
	logger.SetPrefix(prefix)
}

func (logger *Logger) setHook() {
	logger.logger.SetOutput(ioutil.Discard)

	logger.logger.AddHook(&WriterHook{
		Writer: logger.out,
		LogLevels: []logrus.Level{
			logrus.PanicLevel,
			logrus.FatalLevel,
			logrus.ErrorLevel,
		},
	})

	logger.logger.AddHook(&WriterHook{
		Writer: os.Stdout,
		LogLevels: []logrus.Level{
			logrus.InfoLevel,
			logrus.DebugLevel,
		},
	})
}

func (logger *Logger) setOut(fName string) {
	file, err := os.OpenFile(fName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	logger.out = os.Stdout
	if err == nil {
		logger.out = file
	} else {
		logger.logger.Errorf("can't set %v file as stdout | Error %v", fName, err)
	}
}

func parseLogLevel(level string) logrus.Level {
	switch level {
	case "info":
		return logrus.InfoLevel
	case "warn":
		return logrus.WarnLevel
	case "error":
		return logrus.ErrorLevel
	default:
		return logrus.InfoLevel
	}
}

func (logger *Logger) Close() error {
	return logger.out.Close()
}

func (logger *Logger) Out() *os.File {
	return logger.out
}

func (logger *Logger) SetPrefix(prefix string) {
	logger.prefix = prefix
}

func (logger *Logger) Infof(format string, args ...interface{}) {
	logger.logger.Infof(fmt.Sprintf("[%v] %v", logger.prefix, format), args)
}

func (logger *Logger) Info(text interface{}) {
	logger.logger.Info(fmt.Sprintf("[%v] %v", logger.prefix, text))
}

func (logger *Logger) Debug(text interface{}) {
	logger.logger.Debug(fmt.Sprintf("[%v] %v", logger.prefix, text))
}

func (logger *Logger) Debugf(text interface{}, args ...interface{}) {
	logger.logger.Debugf(fmt.Sprintf("[%v] %v", logger.prefix, text), args)
}

func (logger *Logger) Error(text interface{}) {
	logger.logger.Info(fmt.Sprintf("[%v] %v", logger.prefix, text))
}

func (logger *Logger) Errorf(format string, args ...interface{}) {
	logger.logger.Infof(fmt.Sprintf("[%v] %v", logger.prefix, format), args)
}
