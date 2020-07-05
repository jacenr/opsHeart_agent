package logger

import (
	"fmt"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	log "github.com/sirupsen/logrus"
	cc "opsHeart_agent/conf"
	"os"
	"strings"
)

var (
	AgentLog = log.New()
	HbsLog   = log.New()
)

var (
	appMode string
	level   log.Level

	// Add new log file here.
	logPathMap = map[string]*log.Logger{
		"/var/log/superops/agent.log": AgentLog,
		"/var/log/superops/hbs.log":   HbsLog,
	}
)

func InitLog() {
	loglevel := cc.GetLogLevel()
	appMode = cc.GetMode()

	switch strings.ToLower(loglevel) {
	case "info":
		level = log.InfoLevel
	case "warning":
		level = log.WarnLevel
	case "error":
		level = log.ErrorLevel
	default:
		level = log.DebugLevel
	}

	for k, v := range logPathMap {
		createLog(k, v)
	}
}

func createLog(path string, l *log.Logger) {
	loggerWriter, err := rotatelogs.New(fmt.Sprintf("%s.%%Y%%m%%d", path), rotatelogs.WithRotationCount(7))
	if err != nil {
		log.Fatal(err.Error())
	}
	if appMode == "dev" {
		l.SetOutput(os.Stdout)
		l.SetFormatter(&log.TextFormatter{FullTimestamp: true, DisableLevelTruncation: true})
	} else {
		l.SetOutput(loggerWriter)
		l.SetFormatter(&log.JSONFormatter{})
	}
	l.SetLevel(level)
	l.SetReportCaller(true)
}
