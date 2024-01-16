package server

import (
	"github.com/gorilla/mux"
	"github.com/h3th-IV/mysticMerch/internal/api"
)

func SetAdminRoutes(router *mux.Router) {
	adminRouter := router.PathPrefix("/admin").Subrouter()

	//routes for admin
	adminRouter.HandleFunc("/dashboard", api.AdminDashboard)

}
