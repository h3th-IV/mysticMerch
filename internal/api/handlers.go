package api

import (
	"fmt"
	"net/http"
)

// home Handler
func Home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Welcome to MysticeMerch")
}

// signUp post form Hadler
func SignUp(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "SignUP was succesfull")
}

// Login Post Handler
func LogIn(w http.ResponseWriter, r *http.Request) {

}

// user add product
func AddProdcut(w http.ResponseWriter, r *http.Request) {

}

// Serch product by query
func SearchProduct(w http.ResponseWriter, r *http.Request) {

}

// veiw product
func ViewProducts(w http.ResponseWriter, r *http.Request) {

}

// update product details
func UpdateProductDetails(w http.ResponseWriter, r *http.Request) {

}

func UserCart(w http.ResponseWriter, r *http.Request) {

}

// edit prduct
func AddtoCart(w http.ResponseWriter, r *http.Request) {

}

// edit prduct
func RemovefromCart(w http.ResponseWriter, r *http.Request) {

}

// admin send mail
func AdminDashboard(w http.ResponseWriter, r *http.Request) {

}
