package utils

import (
    "log"
    "os"
)

// Logger структура для ведения логов
type Logger struct {
    infoLog  *log.Logger
    errorLog *log.Logger
}

// NewLogger создает новый логгер
func NewLogger() *Logger {
    infoLog := log.New(os.Stdout, "[INFO] ", log.Ldate|log.Ltime|log.Lshortfile)
    errorLog := log.New(os.Stderr, "[ERROR] ", log.Ldate|log.Ltime|log.Lshortfile)

    return &Logger{
        infoLog:  infoLog,
        errorLog: errorLog,
    }
}

// Info пишет сообщение в лог-файл с уровнем INFO
func (l *Logger) Info(msg string) {
    l.infoLog.Println(msg)
}

// Error пишет сообщение в лог-файл с уровнем ERROR
func (l *Logger) Error(msg string) {
    l.errorLog.Println(msg)
}