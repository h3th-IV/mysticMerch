package api

import (
	"database/sql"
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
	defer dataBase.CloseDB()
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
		return
	}

	firstName := r.FormValue("first_name")
	lastName := r.FormValue("last_name")
	email := r.FormValue("email")
	passowrd := r.FormValue("password")
	phoneNumber := r.FormValue("phone_number")

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
		return
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
	defer dataBase.CloseDB()
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
		http.Error(w, "Failed to encode json object: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

//Cart Operations

// view user cart ##
func UserCart(w http.ResponseWriter, r *http.Request) {
	defer dataBase.CloseDB()
	uuid := (r.Context().Value(utils.UserIDkey)).(string)
	user, err := dataBase.GetUserbyUUID(uuid)
	if err != nil {
		http.Error(w, "err getting user_id", http.StatusBadRequest)
		return
	}

	Products, err := dataBase.GetUserCart(*user.ID)
	if err != nil {
		http.Error(w, "Unable to get user's cart", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(Products); err != nil {
		http.Error(w, "Failed to encode json object: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

// edit prduct ##
func AddtoCart(w http.ResponseWriter, r *http.Request) {
	defer dataBase.CloseDB()

	var product *models.RequestProduct
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		http.Error(w, "Failed to decode request: "+err.Error(), http.StatusInternalServerError)
		return
	}
	//get user Id from token
	uuid := r.Context().Value(utils.UserIDkey).(string)
	user, err := dataBase.GetUserbyUUID(uuid)
	if err != nil {
		http.Error(w, "User Possibly Not Authenticated", http.StatusUnauthorized)
		return
	}

	err = dataBase.AddProductoCart(*user.ID, product.Quantity, product.ProductUUID, product.Color, product.Size)
	if err != nil {
		http.Error(w, "Failed to add Product to user cart", http.StatusInternalServerError)
		return
	}

	response := make(map[string]interface{})
	response["message"] = "Product added to user cart succesfully"
	response["product"] = product

	//set content Type#
	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode json object: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

// update cart details like add quantity, change color and size
func UpdateProductDetails(w http.ResponseWriter, r *http.Request) {
	defer dataBase.CloseDB()
	//parse Update details
	var updateDetails *models.RequestProduct
	if err := json.NewDecoder(r.Body).Decode(&updateDetails); err != nil {
		http.Error(w, "Failed to decode object: "+err.Error(), http.StatusBadRequest)
		return
	}

	//get user id from contxt
	uuid := r.Context().Value(utils.UserIDkey).(string)
	user, err := dataBase.GetUserbyUUID(uuid)
	if err != nil {
		http.Error(w, "Failed to retreive user ID: "+err.Error(), http.StatusInternalServerError)
		return
	}

	//get product
	product, err := dataBase.GetProduct(updateDetails.ProductUUID)
	//check if product exist in user cart
	exist, err := dataBase.CheckProductExistInUserCart(*user.ID, *product.ID)
	if err != nil {
		http.Error(w, "Failed to check if Product exist in user's cart"+err.Error(), http.StatusInternalServerError)
		return
	}
	if !exist {
		http.Error(w, "Product not found in user's cart"+err.Error(), http.StatusInternalServerError)
	}

	//update Product details
	if err = dataBase.EditCartItem(*user.ID, *product.ID, updateDetails.Quantity, updateDetails.Color, updateDetails.Size); err != nil {
		http.Error(w, "Failed to update product in user's cart: "+err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"message": "Product details updated succesfully",
	}
	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode json object"+err.Error(), http.StatusInternalServerError)
		return
	}
}

// edit prduct ##
func RemovefromCart(w http.ResponseWriter, r *http.Request) {
	defer dataBase.CloseDB()
	var product *models.RequestProduct
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		http.Error(w, "Failed to decode object"+err.Error(), http.StatusBadRequest)
		return
	}
	uuid := r.Context().Value(utils.UserIDkey).(string)
	user, err := dataBase.GetUserbyUUID(uuid)
	if err != nil {
		http.Error(w, "Failed to get user ID"+err.Error(), http.StatusInternalServerError)
		return
	}
	cartProduct, err := dataBase.GetProduct(product.ProductUUID)
	if err != nil {
		http.Error(w, "Failed to get product from store: "+err.Error(), http.StatusInternalServerError)
	}
	exist, err := dataBase.CheckProductExistInUserCart(*user.ID, *cartProduct.ID)
	if err != nil {
		http.Error(w, "Failed to check if Product exist in user's cart"+err.Error(), http.StatusInternalServerError)
		return
	}
	//check existence of product
	if !exist {
		http.Error(w, "Product not found in user's cart"+err.Error(), http.StatusInternalServerError)
		return
	}

	//check if product is a store item
	cartItem, err := dataBase.GetProduct(product.ProductUUID)
	if err := dataBase.RemoveItemfromCart(*user.ID, *cartItem.ID); err != nil {
		http.Error(w, "Failed to remove item from cart"+err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"message": "Item reomved from cart successfully",
	}
	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode json object: "+err.Error(), http.StatusInternalServerError)
		return
	}

}

// GetItem from cart use list UserPorduct here ##
func GetItemFromCart(w http.ResponseWriter, r *http.Request) {
	dataBase.CloseDB()
	var product *models.RequestProduct
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		http.Error(w, "Failed to decode object: "+err.Error(), http.StatusBadRequest)
		return
	}

	uuid := r.Context().Value(utils.UserIDkey).(string)
	user, err := dataBase.GetUserbyUUID(uuid)
	if err != nil {
		http.Error(w, "Failed to get user ID: "+err.Error(), http.StatusInternalServerError)
		return
	}
	dbPoduct, err := dataBase.GetProduct(product.ProductUUID)
	item, err := dataBase.GetItemFromCart(*user.ID, *dbPoduct.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Item not found in user's cart", http.StatusInternalServerError)
			return
		}
		http.Error(w, "Failed to get item from user's cart: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(item); err != nil {
		http.Error(w, "Failed to encode item ito json object: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

// add new address for user
func AddNewAddr(w http.ResponseWriter, r *http.Request) {
	defer dataBase.CloseDB()
	uuid := r.Context().Value(utils.UserIDkey).(string)

	user, err := dataBase.GetUserbyUUID(uuid)
	if err != nil {
		http.Error(w, "Failed to retrieve user ID", http.StatusInternalServerError)
		return
	}
	if err := r.ParseForm(); err != nil {
		utils.ServerError(w, err)
		return
	}
	house_no := r.FormValue("house_no")
	street := r.FormValue("street")
	city := r.FormValue("city")
	postal_code := r.FormValue("postal_code")

	if err = dataBase.AddUserAddress(user, house_no, street, city, postal_code); err != nil {
		http.Error(w, "Failed to Add new asddress", http.StatusInternalServerError)
		return
	}
	response := map[string]interface{}{
		"message": "Address succefully added",
	}

	w.Header().Set("Content-Tyep", "application/json")
	if err = json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response: "+err.Error(), http.StatusInternalServerError)
		return
	}
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
