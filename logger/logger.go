package logger

import (
	"fmt"
	"github.com/sadlil/gologger"
)

var (
	Logger gologger.GoLogger
)

func init() {
	Logger = gologger.GetLogger(gologger.CONSOLE, gologger.ColoredLog)
}

func SetLogger(newLogger gologger.GoLogger) {
	Logger = newLogger
}

func HandleError(err error, msg string) {
	if err != nil {
		lg := fmt.Sprintf("%s - %s", msg, err)
		Logger.Error(lg)
	}
}
