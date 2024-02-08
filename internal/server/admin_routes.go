package server

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/h3th-IV/mysticMerch/internal/api"
	"github.com/h3th-IV/mysticMerch/internal/utils"
	"github.com/justinas/alice"
)

func SetAdminRoutes(router *mux.Router) {
	adminRouter := router.PathPrefix("/admin").Subrouter()

	//routes for admin
	authChain := alice.New(utils.AdminRoute)
	adminRouter.Handle("/broadcast", authChain.ThenFunc(api.AdminBroadcast)).Methods(http.MethodPost)
	adminRouter.Handle("/transactional", authChain.ThenFunc(api.Transactional)).Methods(http.MethodPost)
}
