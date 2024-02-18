package server

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/h3th-IV/mysticMerch/internal/api"
	"github.com/h3th-IV/mysticMerch/internal/utils"
	"github.com/justinas/alice"
	"go.uber.org/zap"
)

func StartServer() {
	// logger := utils.NewLogger(os.Stdout, os.Stderr)
	//use alice to package potential middleware
	if err := utils.LoadEnv(); err != nil {
		return
	}
	middlewareChain := alice.New(utils.RequestLogger, utils.RecoverPanic)

	router := mux.NewRouter()
	router.HandleFunc("/", api.Home)
	router.HandleFunc("/signup", api.SignUp).Methods(http.MethodPost)
	router.HandleFunc("/login", api.LogIn).Methods(http.MethodPost)

	//set Admin related routes
	SetAdminRoutes(router)

	//set User Related routes
	SetUserRoutes(router)

	//set Product related routes
	SetProductRoutes(router)

	//set Cart routes
	SetCartRoutes(router)

	router.Use(middlewareChain.Then)
	server := &http.Server{
		Addr:     ":8000",
		Handler:  router,
		ErrorLog: zap.NewStdLog(utils.ReplaceLogger),
	}
	utils.ReplaceLogger.Info("Listening and serving :8000")
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		utils.ReplaceLogger.Fatal("Server Failed to start", zap.Error(err))
	}
}
