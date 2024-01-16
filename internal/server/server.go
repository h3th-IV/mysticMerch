package server

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

// struct for the Apllication related configuration
type MarketPlace struct {
	infolog *log.Logger
	errlog  *log.Logger
}

func Routes() {
	//logger for wrting informational message
	InfoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)

	//logger for writing error related messages
	ErrorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	//new merch commerce
	Merch := MarketPlace{
		infolog: InfoLog,
		errlog:  ErrorLog,
	}

	router := mux.NewRouter()

	//set Admin related routes
	SetAdminRoutes(router)

	//set User Related routes
	SetUserRoutes(router)

	//set Product realted routes
	SetProductRoutes(router)

	server := &http.Server{
		Addr:     "8000",
		Handler:  router,
		ErrorLog: Merch.errlog,
	}

	InfoLog.Println("Listening and serving :8000")
	ErrorLog.Fatal(server.ListenAndServe())
}
