package server

import (
	"github.com/gorilla/mux"
	"github.com/h3th-IV/mysticMerch/internal/api"
)

func SetProductRoutes(router *mux.Router) {
	ProductRoutes := router.PathPrefix("/products").Subrouter()

	ProductRoutes.HandleFunc("/viewall", api.ViewProducts)
	ProductRoutes.HandleFunc("/search/{query:[tbb]}", api.SearchProduct)
	ProductRoutes.HandleFunc("/cart", api.UserCart)     //will require authentication MW
	ProductRoutes.HandleFunc("/additem", api.AddtoCart) //will require authentication MW
	ProductRoutes.HandleFunc("/removeitem", api.RemovefromCart)

}
