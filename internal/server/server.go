package server

import (
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/h3th-IV/mysticMerch/internal/api"
	"github.com/h3th-IV/mysticMerch/internal/models"
	"github.com/h3th-IV/mysticMerch/internal/utils"
	"github.com/justinas/alice"
)

// struct for the Apllication related configuration
type MarketPlace struct {
	databox *models.MarketModel
}

func Routes() {
	logger := utils.NewLogger(os.Stdout, os.Stderr)
	//use alice to package potential middleware
	middlewareChain := alice.New(utils.RequestLogger)

	router := mux.NewRouter()
	router.HandleFunc("/", api.Home)
	router.HandleFunc("/user/signup", api.SignUp)

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
		ErrorLog: logger.ErrLogger,
	}
	logger.InfoLogger.Println("Listening and serving :8000")
	logger.ErrLogger.Fatal(server.ListenAndServe())
}
