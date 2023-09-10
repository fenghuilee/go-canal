package log

import (
	"github.com/sirupsen/logrus"
	"os"
)

const (
	LevelPanic logrus.Level = iota
	LevelFatal
	LevelError
	LevelWarn
	LevelInfo
	LevelDebug
	LevelTrace
)

type TConfig struct {
	Level string `yaml:"level"`
}

var Config = &struct {
	Logger *TConfig `yaml:"logger"`
}{
	Logger: new(TConfig),
}

var Logger = logrus.New()

func init() {
	Logger.Out = os.Stdout
	formatter := Logger.Formatter.(*logrus.TextFormatter)
	formatter.ForceColors = true
	formatter.FullTimestamp = true
	formatter.TimestampFormat = "2006-01-02 15:04:05"
}
func SetLevel(l string) {
	level, _ := logrus.ParseLevel(l)
	if &level == nil {
		level = LevelInfo
	}
	Logger.SetLevel(level)
}

func Print(args ...interface{}) {
	Logger.Print(args...)
}

func Printf(format string, args ...interface{}) {
	Logger.Printf(format, args...)
}

func Debug(args ...interface{}) {
	Logger.Debug(args)
}

func Debugf(format string, args ...interface{}) {
	Logger.Debugf(format, args...)
}

func Info(args ...interface{}) {
	Logger.Info(args)
}

func Infof(format string, args ...interface{}) {
	Logger.Infof(format, args...)
}

func Warn(args ...interface{}) {
	Logger.Warn(args)
}

func Warnf(format string, args ...interface{}) {
	Logger.Warnf(format, args...)
}

func Error(args ...interface{}) {
	Logger.Error(args)
}

func Errorf(format string, args ...interface{}) {
	Logger.Errorf(format, args...)
}

func Fatal(args ...interface{}) {
	Logger.Fatal(args)
}

func Fatalf(format string, args ...interface{}) {
	Logger.Fatalf(format, args...)
}

func Panic(args ...interface{}) {
	Logger.Panic(args)
}

func Panicf(format string, args ...interface{}) {
	Logger.Panicf(format, args...)
}
