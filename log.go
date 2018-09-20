package toolkit

import (

	//"io"
	"errors"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type logItem struct {
	LogType string
	Msg     string
}

type LogEngine struct {
	LogToStdOut     bool
	LogToFile       bool
	Path            string
	FileNamePattern string
	UseDateFormat   string

	logInfo  *log.Logger
	logWarn  *log.Logger
	logError *log.Logger

	chanLogItem chan logItem
	fileNames   map[string]string
	writers     map[string]*os.File
	hooks       map[string][]func(string, string)

	//logFile         *log.Logger
	//logFileHandler  *os.File
}

func NewLog(toStdOut bool, toFile bool, path string, fileNamePattern string, useDateFormate string) (*LogEngine, error) {
	var e error
	l := new(LogEngine)
	l.LogToStdOut = toStdOut
	l.LogToFile = toFile
	l.Path = path
	l.FileNamePattern = fileNamePattern
	l.UseDateFormat = useDateFormate
	//l.logger = log.New(out, prefix, flag)

	e = l.initLogger()
	if e != nil {
		return nil, e
	}

	if l.LogToFile {
		l.chanLogItem = make(chan logItem)

		go func() {
			for li := range l.chanLogItem {
				l.writeLogToFile(li.Msg, li.LogType)
			}
		}()
	}

	return l, nil
}

func (l *LogEngine) initLogger() error {
	//var e error = nil
	if l.LogToStdOut {
		l.logError = log.New(os.Stdout, "ERROR ", log.Ldate|log.Ltime)
		l.logInfo = log.New(os.Stdout, "INFO ", log.Ldate|log.Ltime)
		l.logWarn = log.New(os.Stdout, "WARNING ", log.Ldate|log.Ltime)
	}

	l.fileNames = map[string]string{}
	l.writers = map[string]*os.File{}
	l.hooks = map[string][]func(string, string){}
	return nil
}

func (l *LogEngine) AddLog(msg string, logtype string) error {
	var e error
	logtype = strings.ToUpper(logtype)

	if l.LogToStdOut {
		if logtype == "ERROR" {
			l.logError.Println(msg)
		} else if logtype == "WARNING" {
			l.logWarn.Println(msg)
		} else {
			l.logInfo.Println(msg)
		}
		if e != nil {
			return errors.New("Log.AddLog Error: " + e.Error())
		}
	}

	if l.LogToFile {
		l.chanLogItem <- logItem{logtype, msg}
	}

	//--- run hook
	go func() {
		hooks := l.hooks[logtype]
		for _, hook := range hooks {
			hook(logtype, msg)
		}
	}()

	return nil
}

func (l *LogEngine) writeLogToFile(msg, logtype string) {
	filename := l.FileNamePattern
	if l.UseDateFormat != "" && strings.Contains(l.FileNamePattern, "$DATE") {
		filename = strings.Replace(l.FileNamePattern, "$DATE", Date2String(time.Now(), l.UseDateFormat), -1)
	}
	if strings.Contains(filename, "$LOGTYPE") {
		filename = strings.Replace(filename, "$LOGTYPE", logtype, -1)
	}
	filename = filepath.Join(l.Path, filename)
	filenameSelected := l.fileNames[logtype]
	if filename != filenameSelected {
		w, exist := l.writers[logtype]
		if exist {
			w.Close()
		}

		f, e := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
		if e != nil {
			return
			//return errors.New("Log.AddLog Error: " + e.Error())
		}
		l.fileNames[logtype] = filename
		l.writers[logtype] = f
	}
	logFile := log.New(l.writers[logtype], logtype+" ", log.Ldate|log.Ltime)
	logFile.Println(msg)
}

func (l *LogEngine) AddHook(fn func(string, string), logtypes ...string) {
	if len(logtypes) == 0 {
		logtypes = []string{"ERROR", "INFO", "WARNING"}
	}

	for _, logtype := range logtypes {
		hooks := l.hooks[logtype]
		hooks = append(hooks, fn)
		l.hooks[logtype] = hooks
	}
}

func (l *LogEngine) Info(msg string) error {
	return l.AddLog(msg, "INFO")
}

func (l *LogEngine) Error(msg string) error {
	return l.AddLog(msg, "ERROR")
}

func (l *LogEngine) Warning(msg string) error {
	return l.AddLog(msg, "WARNING")
}

func (l *LogEngine) Infof(msg string, args ...interface{}) error {
	msg = Sprintf(msg, args...)
	return l.AddLog(msg, "INFO")
}

func (l *LogEngine) Errorf(msg string, args ...interface{}) error {
	msg = Sprintf(msg, args...)
	return l.AddLog(msg, "ERROR")
}

func (l *LogEngine) Warningf(msg string, args ...interface{}) error {
	msg = Sprintf(msg, args...)
	return l.AddLog(msg, "WARNING")
}

func (l *LogEngine) Close() {
	//l.logFileHandler.Close()
	for _, w := range l.writers {
		w.Close()
	}

	if l.chanLogItem != nil {
		close(l.chanLogItem)
	}
}
