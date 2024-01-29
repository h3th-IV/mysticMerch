package server

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/h3th-IV/mysticMerch/internal/api"
)

func SetProductRoutes(router *mux.Router) {
	ProductRoutes := router.PathPrefix("/products").Subrouter()

	ProductRoutes.HandleFunc("/{id:[a-zA-Z0-9-]+}", api.ViewProduct)
	ProductRoutes.HandleFunc("/search", api.SearchProduct).Methods(http.MethodGet).Queries("product_name")

}

func SetCartRoutes(router *mux.Router) {
	CartProdcts := router.PathPrefix("/carts").Subrouter()

	//cart operations //will require authentication MW
	CartProdcts.HandleFunc("/cart", api.UserCart)
	CartProdcts.HandleFunc("/additem", api.AddtoCart)
	CartProdcts.HandleFunc("/removeitem", api.RemovefromCart)
	CartProdcts.HandleFunc("/itemdetails", api.UpdateProductDetails)
	CartProdcts.HandleFunc("/item", api.GetItemFromCart) //with request query
	CartProdcts.HandleFunc("checkout", api.BuyFromCart)
	CartProdcts.HandleFunc("/buy", api.InstantBuy)
}
