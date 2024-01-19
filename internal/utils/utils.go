package utils

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"golang.org/x/crypto/bcrypt"
)

// Logger represents a custom logger.
type Logger struct {
	InfoLogger *log.Logger
	ErrLogger  *log.Logger
}

// NewLogger creates a new instance of Logger with customizable log formatting.
func NewLogger(infoOutput, errorOutput *os.File) *Logger {
	return &Logger{
		InfoLogger: log.New(infoOutput, "INFO\t", log.Ldate|log.Ltime),
		ErrLogger:  log.New(errorOutput, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile),
	}
}

// LogInfo logs an informational message.
func (l *Logger) LogInfo(message string) {
	l.InfoLogger.Println(message)
}

// LogError logs an error message.
func (l *Logger) LogError(message string) {
	l.ErrLogger.Println(message)
}

// Middleware to log requests to the server
func RequestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger := NewLogger(os.Stdout, os.Stderr)
		logger.LogInfo(fmt.Sprintf("%v - %v %v %v", r.RemoteAddr, r.Proto, r.Method, r.URL.RequestURI()))
		next.ServeHTTP(w, r)
	})
}

// func HashPassword
func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return "", nil
	}
	return string(hash), nil
}
