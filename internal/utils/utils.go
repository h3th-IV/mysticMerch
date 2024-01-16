package utils

import (
	"log"
	"net/http"
	"os"
)

var (
	InfoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
)

// middleware to log request to server
func ReuestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		InfoLog.Printf("%v -%v %v %v", r.RemoteAddr, r.Proto, r.Method, r.URL.RequestURI())
		next.ServeHTTP(w, r)
	})
}
