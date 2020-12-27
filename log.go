package toolkit

import (

	//"io"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type LogItem struct {
	LogType string
	Msg     string
}

type LevelBit int

const (
	AllLevel     int = 0
	InfoLevel        = 1
	WarningLevel     = 2
	ErrorLevel       = 3
	DebugLevel       = 4
)

type LogEngine struct {
	LogToStdOut     bool
	LogToFile       bool
	Path            string
	FileNamePattern string
	UseDateFormat   string

	logInfo  *log.Logger
	logWarn  *log.Logger
	logError *log.Logger
	logDebug *log.Logger

	chanLogItem chan LogItem
	fileNames   map[string]string
	writers     map[string]*os.File
	hooks       map[string][]func(string, string)

	stdOutLevels  []bool
	fileOutLevels []bool

	prefix     string
	fnTemplate func(LogItem) string

	//logFile         *log.Logger
	//logFileHandler  *os.File
}

//type LogFields map[string]interface{}

func NewLog(toStdOut bool, toFile bool, path string, fileNamePattern string, useDateFormat string) (*LogEngine, error) {
	var e error
	l := new(LogEngine)
	l.LogToStdOut = toStdOut
	l.LogToFile = toFile
	l.Path = path
	l.FileNamePattern = fileNamePattern
	l.UseDateFormat = useDateFormat
	//l.logger = log.New(out, prefix, flag)

	l.stdOutLevels = make([]bool, 5)
	l.fileOutLevels = make([]bool, 5)

	l.SetLevelStdOuts(InfoLevel, WarningLevel, ErrorLevel)
	l.SetLevelFiles(InfoLevel, WarningLevel, ErrorLevel)

	e = l.initLogger()
	if e != nil {
		return nil, e
	}

	if l.LogToFile {
		l.chanLogItem = make(chan LogItem)

		go func() {
			for li := range l.chanLogItem {
				l.writeLogToFile(li.Msg, li.LogType)
			}
		}()
	}

	return l, nil
}

func NewLogEngine(toStdOut bool, toFile bool, path string, fileNamePattern string, useDateFormat string) *LogEngine {
	l, _ := NewLog(toStdOut, toFile, path, fileNamePattern, useDateFormat)
	return l
}

func (l *LogEngine) initLogger() error {
	//var e error = nil
	l.initStdOut()
	l.fileNames = map[string]string{}
	l.writers = map[string]*os.File{}
	l.hooks = map[string][]func(string, string){}
	return nil
}

func (l *LogEngine) initStdOut() {
	if l.LogToStdOut {
		l.logError = prepareStdoutLogger(l, "ERROR")
		l.logInfo = prepareStdoutLogger(l, "INFO")
		l.logWarn = prepareStdoutLogger(l, "WARNING")
		l.logDebug = prepareStdoutLogger(l, "DEBUG")
	}
}

func prepareStdoutLogger(l *LogEngine, logType string) *log.Logger {
	logger := new(log.Logger)
	logger.SetFlags(0)

	w := new(LogWriter)
	w.initialItem = LogItem{}
	w.initialItem.LogType = logType
	msg := l.prefix
	if l.fnTemplate == nil {
		if msg != "" {
			msg += " "
		}
		msg += "{TIME} "
		msg += logType
		msg += " {MSG}"

		w.fn = func(item LogItem) string {
			fmtMsg := msg
			fmtMsg = strings.Replace(fmtMsg, "{TIME}", time.Now().Format(time.RFC3339), -1)
			fmtMsg = strings.Replace(fmtMsg, "{MSG}", item.Msg, -1)
			return fmtMsg
		}
	} else {
		w.fn = l.fnTemplate
	}

	logger.SetOutput(w)
	return logger
}

func (l *LogEngine) SetStdoutTemplate(fnTemplate func(LogItem) string) {
	l.fnTemplate = fnTemplate
	l.initStdOut()
}

func (l *LogEngine) HasTemplate() bool {
	return l.fnTemplate != nil
}

func (l *LogEngine) SetPrefix(s string) *LogEngine {
	l.prefix = s
	l.initStdOut()
	return l
}

func (l *LogEngine) Prefix() string {
	return l.prefix
}

func (l *LogEngine) SetLevelStdOuts(levels ...int) *LogEngine {
	for i := range []int{0, 1, 2, 3, 4} {
		l.stdOutLevels[i] = false
	}

	for _, level := range levels {
		if level != AllLevel {
			l.stdOutLevels[AllLevel] = false
		}
		l.stdOutLevels[level] = true
	}

	return l
}

func (l *LogEngine) SetLevelStdOut(level int, value bool) *LogEngine {
	if level != AllLevel {
		l.stdOutLevels[AllLevel] = false
	}
	l.stdOutLevels[level] = value
	return l
}

func (l *LogEngine) SetLevelFiles(levels ...int) *LogEngine {
	for i := range []int{0, 1, 2, 3, 4} {
		l.fileOutLevels[i] = false
	}

	for _, level := range levels {
		if level != AllLevel {
			l.fileOutLevels[AllLevel] = false
		}
		l.fileOutLevels[level] = true
	}

	return l
}

func (l *LogEngine) SetLevelFile(level int, value bool) *LogEngine {
	if level != AllLevel {
		l.fileOutLevels[AllLevel] = false
	}
	l.fileOutLevels[level] = value

	return l
}

func (l *LogEngine) StdOutLevel(level int) bool {
	return l.stdOutLevels[level]
}

func (l *LogEngine) FileOutLevel(level int) bool {
	return l.fileOutLevels[level]
}

func (l *LogEngine) AddLog(msg string, logtype string) error {
	var e error
	logtype = strings.ToUpper(logtype)

	if l.LogToStdOut {
		if logtype == "ERROR" && (l.StdOutLevel(AllLevel) || l.StdOutLevel(ErrorLevel)) {
			l.logError.Println(msg)
		} else if logtype == "WARNING" && (l.StdOutLevel(AllLevel) || l.StdOutLevel(WarningLevel)) {
			l.logWarn.Println(msg)
		} else if logtype == "DEBUG" && (l.StdOutLevel(AllLevel) || l.StdOutLevel(DebugLevel)) {
			l.logDebug.Println(msg)
		} else if logtype == "INFO" && (l.StdOutLevel(AllLevel) || l.StdOutLevel(InfoLevel)) {
			l.logInfo.Println(msg)
		}
		if e != nil {
			return errors.New("Log.AddLog Error: " + e.Error())
		}
	}

	if l.LogToFile {
		l.chanLogItem <- LogItem{logtype, msg}
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
	if logtype == "ERROR" && !l.FileOutLevel(AllLevel) && !l.FileOutLevel(ErrorLevel) {
		return
	} else if logtype == "WARNING" && !l.FileOutLevel(AllLevel) && !l.FileOutLevel(WarningLevel) {
		return
	} else if logtype == "INFO" && !l.FileOutLevel(AllLevel) && !l.FileOutLevel(InfoLevel) {
		return
	} else if logtype == "DEBUG" && !l.FileOutLevel(AllLevel) && !l.FileOutLevel(DebugLevel) {
		return
	}

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

	if l.prefix == "" {
		logFile := log.New(l.writers[logtype], logtype+" ", log.Ldate|log.Ltime)
		logFile.Println(msg)
	} else {
		logFile := log.New(l.writers[logtype], l.prefix+" "+logtype+" ", log.Ldate|log.Ltime)
		logFile.Println(msg)
	}
}

func (l *LogEngine) AddHook(fn func(string, string), logtypes ...string) {
	if len(logtypes) == 0 {
		logtypes = []string{"ERROR", "INFO", "WARNING", "DEBUG"}
	}

	for _, logtype := range logtypes {
		hooks := l.hooks[logtype]
		hooks = append(hooks, fn)
		l.hooks[logtype] = hooks
	}
}

func (l *LogEngine) Debug(msg string) error {
	return l.AddLog(msg, "DEBUG")
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

func (l *LogEngine) Debugf(msg string, args ...interface{}) error {
	msg = Sprintf(msg, args...)
	return l.AddLog(msg, "DEBUG")
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

// Error2 will send msg1 to output log and msg2 to system, it could be useful for logging something that can only be seen by sysadmin and user
func (l *LogEngine) Error2(msg1, msg2 string, parm ...interface{}) error {
	l.Errorf(msg2, parm...)
	return errors.New(msg1)
}

/* deprecated
func LogM(m M, msg string) string {
	return Sprintf("field:%s message:%s",
		JsonString(m), msg)
}
*/

var _logger *LogEngine

func Logger() *LogEngine {
	if _logger == nil {
		_logger, _ = NewLog(true, false, "", "", "")
	}
	return _logger
}

type LogWriter struct {
	initialItem LogItem
	fn          func(item LogItem) string
}

func (w *LogWriter) Write(bs []byte) (int, error) {
	item := w.initialItem
	item.Msg = string(bs)
	return fmt.Print(w.fn(item))
}
