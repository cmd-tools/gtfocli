package logger

import (
	"fmt"
	"github.com/sirupsen/logrus"
	easy "github.com/t-tomalak/logrus-easy-formatter"
	"io"
	"os"
	"path"
	"time"
)

var Logger = logrus.New()

func Init(commandName string, debug bool) {

	if !debug {
		Logger.SetOutput(io.Discard)
		return
	}

	logFile := "gtfocli-" + time.Now().Format("20060102") + ".log"
	logLevel := logrus.InfoLevel

	basePath := path.Dir(logFile)
	if err := os.MkdirAll(basePath, 0777); err != nil {
	}

	var f *os.File
	var err error

	if f, err = os.OpenFile(logFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666); err != nil {
		fmt.Println(err)
	}

	Logger.SetOutput(f)

	Logger.SetLevel(logLevel)

	// Setup the logger format
	formatter := &easy.Formatter{
		TimestampFormat: "2006-01-02 15:04:05",
		LogFormat:       "[%lvl%]: %time% - Command: " + commandName + " - %msg%\n",
	}

	Logger.SetFormatter(formatter)
}
