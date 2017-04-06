package gotable

import (
	"fmt"
	"io"
	"log"
)

// level precedence
const (
	debugLevel = 1
	infoLevel  = 2
	errorLevel = 3
)

// level map
var levelInfo = map[string]int{
	"debug": debugLevel,
	"info":  infoLevel,
	"error": errorLevel,
}

// gotable log struct
var gtLog struct {
	logger *log.Logger
	level  int
}

// SetLogger must be called first before using any method of gotable logging
func SetLogger(w io.Writer, level string) {

	// take error as an default level
	if lvl, ok := levelInfo[level]; ok {
		gtLog.level = lvl
	} else {
		gtLog.level = levelInfo["error"]
	}

	gtLog.logger = log.New(w, "GOTABLE ", log.Ldate|log.Ltime|log.Lshortfile)
}

// error message in the log
func errorLog(msg ...interface{}) {

	// first verify the condition, if it does not match then simply return
	if !(gtLog.level > 0 && gtLog.logger != nil) {
		return
	}

	// print messages in logger
	gtLog.logger.Output(2, "ERROR: "+fmt.Sprint(msg...))
}

// info message in the log
func infoLog(msg ...interface{}) {

	// first verify the condition, if it does not match then simply return
	if !(gtLog.level <= infoLevel && gtLog.level > 0 && gtLog.logger != nil) {
		return
	}

	// printing info messages
	gtLog.logger.Output(2, "INFO: "+fmt.Sprint(msg...))
}

// debug message in the log
func debugLog(msg ...interface{}) {

	// first verify the condition, if it does not match then simply return
	if !(gtLog.level <= debugLevel && gtLog.level > 0 && gtLog.logger != nil) {
		return
	}

	// print debugging messages in logger
	gtLog.logger.Output(2, "DEBUG: "+fmt.Sprint(msg...))
}
