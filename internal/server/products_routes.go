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
	CartProdcts := router.PathPrefix("/carts").Subrouter()

	userMWchain := alice.New(utils.AuthRoute)
	//cart operations //will require authentication MW
	CartProdcts.Handle("/cart", userMWchain.ThenFunc(api.GetUserCart)).Methods(http.MethodGet)
	CartProdcts.Handle("/additem", userMWchain.ThenFunc(api.AddtoCart)).Methods(http.MethodPost)
	CartProdcts.Handle("/updateitem", userMWchain.ThenFunc(api.UpdateProductDetails)).Methods(http.MethodPut)
	CartProdcts.Handle("/removeitem", userMWchain.ThenFunc(api.RemovefromCart)).Methods(http.MethodDelete)
	CartProdcts.Handle("/item", userMWchain.ThenFunc(api.GetItemFromCart)).Methods(http.MethodGet)
	CartProdcts.Handle("checkout", userMWchain.ThenFunc(api.BuyFromCart))
	CartProdcts.Handle("/buy", userMWchain.ThenFunc(api.InstantBuy))
}
