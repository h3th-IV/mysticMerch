package utils

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime/debug"

	"github.com/go-sql-driver/mysql"
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

var (
	ErrNoRecord = errors.New("err: no matching record found")

	ErrInvalidCredentials = errors.New("err: invalid credentials")

	ErrExsistingCrednetials = errors.New("err: dupliacte Credentials")
	MySQLErr                *mysql.MySQLError
)

// recover panic
func RecoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				//if panic close connection
				w.Header().Set("Connection", "Close")
				//write internal server error
				ServerError(w, fmt.Errorf("%v", err))
			}
		}()
		next.ServeHTTP(w, r)
	})
}

// used for all internal server Error
func ServerError(w http.ResponseWriter, err error) {
	logger := NewLogger(os.Stdout, os.Stderr)
	errTrace := fmt.Sprintf("%v\n%v", err.Error(), debug.Stack())
	//write output for logging event 2 step backwards
	logger.ErrLogger.Output(2, errTrace)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}
