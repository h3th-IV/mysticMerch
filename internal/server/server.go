package server

import (
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/h3th-IV/mysticMerch/internal/api"
	"github.com/h3th-IV/mysticMerch/internal/utils"
	"github.com/justinas/alice"
)

// struct for the Apllication related configuration
type MarketPlace struct {
}

func Routes() {
	logger := utils.NewLogger(os.Stdout, os.Stderr)
	//use alice to package potential middleware
	middlewareChain := alice.New(utils.RequestLogger, utils.RecoverPanic)

	router := mux.NewRouter()
	router.HandleFunc("/", api.Home)

	//set Admin related routes
	SetAdminRoutes(router)

	//set User Related routes
	SetUserRoutes(router)

	//set Product realted routes
	SetProductRoutes(router)

	//set Cart routes
	SetCartRoutes(router)

	router.Use(middlewareChain.Then)
	server := &http.Server{
		Addr:     ":8000",
		Handler:  router,
		ErrorLog: logger.ErrLogger,
	}
	logger.InfoLogger.Println("Listening and serving :8000")
	logger.ErrLogger.Fatal(server.ListenAndServe())
}
