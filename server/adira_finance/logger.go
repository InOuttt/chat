package adira_finance

import (
	"log"
	"path/filepath"
	"runtime"

	"github.com/sirupsen/logrus"
)

const (
	// InfoColor ...
	InfoColor = "\033[1;34m%s\033[0m" // blue
	// NoticeColor ...
	NoticeColor = "\033[1;36m%s\033[0m"
	// WarningColor ...
	WarningColor = "\033[1;33m%s\033[0m"
	// ErrorColor ...
	ErrorColor = "\033[1;31m%s\033[0m" // red
	// DebugColor ...
	DebugColor = "\033[0;36m%s\033[0m"
)

var (
	// Log ...
	Log = LogS{
		Logger: logrus.New(),
	}
	LogInfo  = Log.Logger.Info
	LogDebug = Log.Logger.Debug
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
	log.Printf(InfoColor, a)
	_, pathfile, no, ok := runtime.Caller(1)
	if ok {
		log.Printf("called by %s#%d", filepath.Base(pathfile), no)
	}
}

// Error ...
// Print error
func (l *LogS) Error(a interface{}) {
	log.Printf(ErrorColor, a)
	_, pathfile, no, ok := runtime.Caller(1)
	if ok {
		log.Printf("called by %s#%d", filepath.Base(pathfile), no)
	}
}

// Warn ...
// Print warn
func (l *LogS) Warn(a interface{}) {
	log.Printf(WarningColor, a)
	_, pathfile, no, ok := runtime.Caller(1)
	if ok {
		log.Printf("called by %s#%d", filepath.Base(pathfile), no)
	}
}

// Notice ...
// Print notice
func (l *LogS) Notice(a interface{}) {
	log.Printf(NoticeColor, a)
	_, pathfile, no, ok := runtime.Caller(1)
	if ok {
		log.Printf("called by %s#%d", filepath.Base(pathfile), no)
	}
}

// Debug ...
// Print debug
func (l *LogS) Debug(a interface{}) {
	log.Printf(DebugColor, a)
	_, pathfile, no, ok := runtime.Caller(1)
	if ok {
		log.Printf("called by %s#%d", filepath.Base(pathfile), no)
	}
}

// Ln ...
// Println
func (l *LogS) Ln(a ...interface{}) {
	log.Println(a...)
}
