package log

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"runtime"
)

var (
	Log *logrus.Logger
)

func init() {
	Log = logrus.New()
	Log.SetOutput(os.Stdout)
	Log.SetFormatter(&logrus.TextFormatter{
		//以下设置只是为了使输出更美观
		DisableColors:   true,
		TimestampFormat: "2006-01-02 15:03:04",
	})
}

func getCallerInfo() (string, int) {
	_, file, line, _ := runtime.Caller(3) // 调用栈层数
	return file, line
}

func getFormat(text string) string {
	file, line := getCallerInfo()
	return fmt.Sprintf("%s position=%s:%d", text, file, line)
}

func Infof(text string, args ...interface{}) {
	format := getFormat(text)
	if args == nil {
		Log.Info(format)
	} else {
		Log.Infof(format, args...)
	}

}

func Info(format string) {
	format = getFormat(format)
	Log.Info(format)
}

func Errorf(format string, args ...interface{}) {
	format = getFormat(format)
	if args == nil {
		Log.Error(format)
	} else {
		Log.Errorf(format, args...)
	}
}

func Error(format string) {
	format = getFormat(format)
	Log.Error(format)
}
