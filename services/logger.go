package services

import (
	"log"
	"os"
)

// Logger define un sistema de logging centralizado
type Logger struct {
	logger *log.Logger
}

// NewLogger crea una nueva instancia del logger
func NewLogger() *Logger {
	return &Logger{
		logger: log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile),
	}
}

// Info loguea un mensaje de información
func (l *Logger) Info(message string) {
	l.logger.SetPrefix("INFO: ")
	l.logger.Println(message)
}

// Warn loguea un mensaje de advertencia
func (l *Logger) Warn(message string) {
	l.logger.SetPrefix("WARN: ")
	l.logger.Println(message)
}

// Error loguea un mensaje de error
func (l *Logger) Error(message string) {
	l.logger.SetPrefix("ERROR: ")
	l.logger.Println(message)
}
