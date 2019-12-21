package adiraFinance

import (
	"log"
	"path/filepath"
	"runtime"

	"github.com/sirupsen/logrus"
)

const (
	// InfoColor ...
	InfoColor = "\033[1;34m%#v\033[0m" // blue
	// NoticeColor ...
	NoticeColor = "\033[1;36m%#v\033[0m" // light blue
	// WarningColor ...
	WarningColor = "\033[1;33m%#v\033[0m" // yellow
	// ErrorColor ...
	ErrorColor = "\033[1;31m%#v\033[0m" // red
	// DebugColor ...
	DebugColor = "\033[0;36m%#v\033[0m" // light blue thin
)

var (
	// Log ...
	Log = LogS{
		Logger: logrus.New(),
	}
	// LogInfo ...
	LogInfo = Log.Logger.Info
	// LogDebug ...
	LogDebug = Log.Logger.Debug
	// LogError ...
	LogError = Log.Logger.Error
)

type (
	// LogS ...
	LogS struct {
		Logger      *logrus.Logger
		totalLogged int64
	}
)

func init() {
	Log.logInit()
}

func (l *LogS) logInit() {
	l.Logger.SetFormatter(&logrus.TextFormatter{
		DisableColors:   false,
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	})
}

// Info ...
// Print info
func (l *LogS) Info(a interface{}) {
	if a == nil {
		a = "[nil]"
	}

	log.Printf(InfoColor, a)
	// _, pathfile, line, ok := runtime.Caller(1)
	// if ok {
	// 	log.Printf("called by %s#%d", filepath.Base(pathfile), line)
	// }
}

// Error ...
// Print error
func (l *LogS) Error(a interface{}) {
	if a == nil {
		a = "[nil]"
	}

	log.Printf(ErrorColor, a)
	_, pathfile, line, ok := runtime.Caller(1)
	if ok {
		log.Printf("called by %s#%d", filepath.Base(pathfile), line)
	}
}

// Warn ...
// Print warn
func (l *LogS) Warn(a interface{}) {
	if a == nil {
		a = "[nil]"
	}

	log.Printf(WarningColor, a)
	_, pathfile, line, ok := runtime.Caller(1)
	if ok {
		log.Printf("called by %s#%d", filepath.Base(pathfile), line)
	}
}

// Notice ...
// Print notice
func (l *LogS) Notice(a interface{}) {
	if a == nil {
		a = "[nil]"
	}

	log.Printf(NoticeColor, a)
	_, pathfile, line, ok := runtime.Caller(1)
	if ok {
		log.Printf("called by %s#%d", filepath.Base(pathfile), line)
	}
}

// Debug ...
// Print debug
func (l *LogS) Debug(a interface{}) {
	if a == nil {
		a = "[nil]"
	}

	log.Printf(DebugColor, a)
	_, pathfile, line, ok := runtime.Caller(1)
	if ok {
		log.Printf("called by %s#%d", filepath.Base(pathfile), line)
	}
}

// Ln ...
// Println
func (l *LogS) Ln(a ...interface{}) {
	log.Println(a...)
}
