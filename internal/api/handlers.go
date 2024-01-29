package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/h3th-IV/mysticMerch/internal/database"
	"github.com/h3th-IV/mysticMerch/internal/models"
	"github.com/h3th-IV/mysticMerch/internal/utils"
)

var (
	db, _    = database.InitDB()
	dataBase = database.DBModel{
		DB: db,
	}
)

// home Handler display a list products
func Home(w http.ResponseWriter, r *http.Request) {
	//get some list of prduct to display on the home page
	products, err := dataBase.ViewProducts()
	if err != nil {
		utils.ServerError(w, err)
	}
	json.NewEncoder(w).Encode(products)
}

// signUp post form Hadler ##
func SignUp(w http.ResponseWriter, r *http.Request) {
	defer dataBase.CloseDB()
	err := r.ParseForm()
	if err != nil {
		utils.ServerError(w, err)
	}

	firstName := r.FormValue("firstName")
	lastName := r.FormValue("lastName")
	email := r.FormValue("email")
	passowrd := r.FormValue("password")
	phoneNumber := r.FormValue("phoneNumber")

	//validate user input as w don't trust user input
	isDetailsValid := utils.ValidateSignUpDetails([]models.ValidAta{
		{Value: firstName, Validator: "fName"}, // "first_name"
		{Value: lastName, Validator: "lName"},
		{Value: email, Validator: "email"},
		{Value: passowrd, Validator: "password"},
	})
	if !isDetailsValid {
		http.Error(w, "Invalid User Input", http.StatusBadRequest)
	}

	err = dataBase.InsertUser(firstName, lastName, email, phoneNumber, passowrd)
	if err != nil {
		utils.ServerError(w, err)
		return
	}
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

// Login Post Handler ##
func LogIn(w http.ResponseWriter, r *http.Request) {
	defer dataBase.CloseDB()
	if err := r.ParseForm(); err != nil {
		utils.ServerError(w, err)
	}
	email := r.FormValue("email")
	password := r.FormValue("password")
	user, err := dataBase.AuthenticateUser(email, password)
	if err != nil {
		if errors.Is(err, utils.ErrInvalidCredentials) {
			http.Error(w, "Invalid email or password", http.StatusUnauthorized)
			return
		}
		utils.ServerError(w, err)
		return
	}
	JWToken, err := utils.GenerateToken(user)
	if err != nil {
		utils.ServerError(w, err)
	}
	//send token to client
	if err := json.NewEncoder(w).Encode(JWToken); err != nil {
		http.Error(w, "Unable to encode token into JSON object"+err.Error(), http.StatusInternalServerError)
	}
}

// Serch product by query name
func SearchProduct(w http.ResponseWriter, r *http.Request) {
	defer dataBase.CloseDB()
	//get name from search query
	query := r.URL.Query()
	ProductName := query.Get("product_name")

	Products, err := dataBase.GetProductByName(ProductName)
	if err != nil {
		http.Error(w, "Product Not Available Yet", http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")

	//encode
	if err = json.NewEncoder(w).Encode(Products); err != nil {
		http.Error(w, "Unable to encode products into JSON:"+err.Error(), http.StatusInternalServerError)
		return
	}
}

// veiw product ##
func ViewProduct(w http.ResponseWriter, r *http.Request) {
	Var := mux.Vars(r)
	Product_id := Var["id"]

	//get by the uuid of product
	Product, err := dataBase.GetProduct(Product_id)
	if err != nil {
		http.Error(w, "Failed to Fetch Product", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(Product); err != nil {
		http.Error(w, "Unable to encode product into JSON:"+err.Error(), http.StatusInternalServerError)
		return
	}
}

//Cart Operations

// view user cart ##
func UserCart(w http.ResponseWriter, r *http.Request) {
	user_id := (r.Context().Value(utils.UserIDkey)).(string)
	id, err := dataBase.GetUserID(user_id)
	if err != nil {
		http.Error(w, "err getting user_id", http.StatusBadRequest)
		return
	}

	Products, err := dataBase.GetUserCart(id)
	if err != nil {
		http.Error(w, "Unable to get user's cart", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(Products); err != nil {
		http.Error(w, "Unable to encode products into JSON:"+err.Error(), http.StatusInternalServerError)
		return
	}
}

// edit prduct ##
func AddtoCart(w http.ResponseWriter, r *http.Request) {

}

// update product details like add quantity
func UpdateProductDetails(w http.ResponseWriter, r *http.Request) {

}

// edit prduct ##
func RemovefromCart(w http.ResponseWriter, r *http.Request) {

}

// GetItem from cart use list UserPorduct here ##
func GetItemFromCart(w http.ResponseWriter, r *http.Request) {

}

// buy from cart ##
func BuyFromCart(w http.ResponseWriter, r *http.Request) {

}

// instant buy ##
func InstantBuy(w http.ResponseWriter, r *http.Request) {

}

// admin send mail ##
func AdminDashboard(w http.ResponseWriter, r *http.Request) {

}
