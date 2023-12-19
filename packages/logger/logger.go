package logger

import (
    "io"
    "os"
    "log"
    "fmt"
)


const (
    LogLevelDebug = iota
    LogLevelInfo
    LogLevelError
)

var currentLogLevel = LogLevelInfo
var logFile *os.File


func InitLogger(logDirPath string) error {
    logFileName := "sf.log"
    logFilePath := logDirPath + "/" + logFileName
    logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
    if err != nil {
        return fmt.Errorf("Failed to open log file at path %s: %v", logFilePath, err)
    }

    multiWriter := io.MultiWriter(os.Stdout, logFile)
    log.SetOutput(multiWriter)
    return nil
}

func CloseLogger() {
    if logFile != nil {
        logFile.Close()
    }
}

func Debug(format string, v ...interface{}) {
    if currentLogLevel <= LogLevelDebug {
        log.Printf("DEBUG: "+format, v...)
    }
}

func Info(format string, v ...interface{}) {
    if currentLogLevel <= LogLevelInfo {
        log.Printf("INFO: "+format, v...)
    }
}

func Error(format string, v ...interface{}) {
    if currentLogLevel <= LogLevelError {
        log.Printf("ERROR: "+format, v...)
    }
}

func Fatal(format string, v ...interface{}) {
    if currentLogLevel <= LogLevelError {
        log.Fatalf("FATAL: "+format, v...)
    }
}
