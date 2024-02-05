package server

import (
	"github.com/gorilla/mux"
	"github.com/h3th-IV/mysticMerch/internal/api"
	"github.com/h3th-IV/mysticMerch/internal/utils"
	"github.com/justinas/alice"
)

func SetAdminRoutes(router *mux.Router) {
	adminRouter := router.PathPrefix("/admin").Subrouter()

	//routes for admin
	authChain := alice.New(utils.AdminRoute)
	adminRouter.Handle("/dashboard", authChain.ThenFunc(api.AdminDashboard))

}
