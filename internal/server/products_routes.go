package server

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/h3th-IV/mysticMerch/internal/api"
	"github.com/h3th-IV/mysticMerch/internal/utils"
	"github.com/justinas/alice"
)

func SetProductRoutes(router *mux.Router) {
	ProductRoutes := router.PathPrefix("/products").Subrouter()

	ProductRoutes.HandleFunc("/{id:[a-zA-Z0-9-]+}", api.ViewProduct).Methods(http.MethodGet)
	ProductRoutes.HandleFunc("/catalog", api.SearchProduct).Methods(http.MethodGet)

}

func SetCartRoutes(router *mux.Router) {
	CartProducts := router.PathPrefix("/carts").Subrouter()

	userMWchain := alice.New(utils.AuthRoute)
	//cart operations //will require authentication MW
	CartProducts.Handle("/cart", userMWchain.ThenFunc(api.GetUserCart)).Methods(http.MethodGet)
	CartProducts.Handle("/additem", userMWchain.ThenFunc(api.AddtoCart)).Methods(http.MethodPost)
	CartProducts.Handle("/updateitem", userMWchain.ThenFunc(api.UpdateProductDetails)).Methods(http.MethodPut)
	CartProducts.Handle("/removeitem", userMWchain.ThenFunc(api.RemovefromCart)).Methods(http.MethodDelete)
	CartProducts.Handle("/item", userMWchain.ThenFunc(api.GetItemFromCart)).Methods(http.MethodGet)
	CartProducts.Handle("checkout", userMWchain.ThenFunc(api.BuyFromCart))
	CartProducts.Handle("/buy", userMWchain.ThenFunc(api.InstantBuy))
}
