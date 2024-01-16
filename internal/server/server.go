package server

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/h3th-IV/mysticMerch/internal/api"
	"github.com/h3th-IV/mysticMerch/internal/utils"
	"github.com/justinas/alice"
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

	//use alice package potential middleware
	middlewareChain := alice.New(utils.ReuestLogger)

	//new merch commerce
	MerchApp := MarketPlace{
		infolog: InfoLog,
		errlog:  ErrorLog,
	}

	router := mux.NewRouter()
	router.HandleFunc("/", api.Home)

	//set Admin related routes
	SetAdminRoutes(router)

	//set User Related routes
	SetUserRoutes(router)

	//set Product realted routes
	SetProductRoutes(router)

	router.Use(middlewareChain.Then)
	server := &http.Server{
		Addr:     ":8000",
		Handler:  router,
		ErrorLog: MerchApp.errlog,
	}
	InfoLog.Println("Listening and serving :8000")
	ErrorLog.Fatal(server.ListenAndServe())
}
