package utils

import (
	"log"
	"os"
	"sync"
	"time"
)

type Logger struct {
	mu     sync.Mutex
	logger *log.Logger
}

var instance *Logger
var once sync.Once // mutex

// GetLogger returns a singleton instance of Logger
func GetLogger() *Logger {
	// once.Do is used to ensure that the instance is created only once
	once.Do(func() {
		instance = &Logger{
			logger: log.New(os.Stdout, "", log.LstdFlags),
		}
	})
	return instance
}

func (l *Logger) Info(format string, v ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.logger.SetPrefix("[INFO] " + time.Now().Format(time.RFC3339))
	l.logger.Printf(format, v...)
}

func (l *Logger) Error(format string, v ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.logger.SetPrefix("[ERROR] " + time.Now().Format(time.RFC3339))
	l.logger.Printf(format, v...)
}