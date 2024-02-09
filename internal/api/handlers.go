package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/h3th-IV/mysticMerch/internal/admin"
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
		utils.ServerError(w, "Failed to get products. Please try again later.", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(products); err != nil {
		utils.ServerError(w, "Failed to encode json object.", err)
		return
	}
}

// admin stuff
func AddItemtoStore(w http.ResponseWriter, r *http.Request) {
	defer dataBase.CloseDB()

	//decode new item
	var Product *models.Product
	if err := json.NewDecoder(r.Body).Decode(&Product); err != nil {
		http.Error(w, "Failed to decode json object", http.StatusBadRequest)
		return
	}

	_, err := dataBase.AddProduct(*Product.ProductName, *Product.Description, *Product.Image, *Product.Price)
	if err != nil {
		utils.ServerError(w, "Failed to add product to store", err)
		return
	}

	response := map[string]interface{}{
		"message": "Operation was succesfull",
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		utils.ServerError(w, "Failed to encode object", err)
		return
	}
}

// admin stuff
func RemoveItemfromStore(w http.ResponseWriter, r *http.Request) {
	defer dataBase.CloseDB()

	//decode item -- out of stock item
	var Product *models.RequestProduct
	if err := json.NewDecoder(r.Body).Decode(&Product); err != nil {
		http.Error(w, "Failed to decode json object"+err.Error(), http.StatusBadRequest)
		return
	}

	if err := dataBase.RemoveProductFromStore(Product.ProductUUID); err != nil {
		utils.ServerError(w, "Failed to reomve itme from store", err)
		return
	}

	response := make(map[string]interface{})
	response["message"] = "Operation was succefull"
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		utils.ServerError(w, "Failed to encode json object", err)
		return
	}
}

// admin send mail ##
func AdminBroadcast(w http.ResponseWriter, r *http.Request) {
	var notification *models.BroadcastNotification
	if err := json.NewDecoder(r.Body).Decode(&notification); err != nil {
		http.Error(w, "Failed to decode json"+err.Error(), http.StatusBadRequest)
		return
	}

	users, err := dataBase.GetAllUsers()
	if err != nil {
		http.Error(w, "Failed to retrive users for Brodcast message"+err.Error(), http.StatusInternalServerError)
		return
	}
	if err := admin.MarketingEmail(users, notification.Subject, notification.Body); err != nil {
		http.Error(w, "Failed to send broadcast message"+err.Error(), http.StatusInternalServerError)
		return
	}
	response := map[string]interface{}{
		"message": "Broadcast email send succesfully",
	}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		utils.ServerError(w, "Failed to encode response", err)
		return
	}
}

// send Transactional email to aparticular Customer
func Transactional(w http.ResponseWriter, r *http.Request) {
	var notification *models.TransactionNotification
	if err := json.NewDecoder(r.Body).Decode(&notification); err != nil {
		http.Error(w, "Failed to decode request"+err.Error(), http.StatusBadRequest)
		return
	}

	if notification.ResponseUser.Email == "" || notification.Body == "" || notification.Subject == "" {
		http.Error(w, "User email, email body or email subject is empty", http.StatusInternalServerError)
		return
	}

	if err := admin.TransactionalEmail(&notification.ResponseUser, notification.Subject, notification.Body); err != nil {
		http.Error(w, "Failed to send email to user"+err.Error(), http.StatusInternalServerError)
		return
	}

	response := make(map[string]interface{})
	response["message"] = "email sent Successfully"

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response"+err.Error(), http.StatusInternalServerError)
		return
	}
}

// signUp post form Hadler ##
func SignUp(w http.ResponseWriter, r *http.Request) {
	defer dataBase.CloseDB()
	err := r.ParseForm()
	if err != nil {
		utils.ServerError(w, "Failed to parse form.", err)
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
		http.Error(w, "Invalid User Input.", http.StatusBadRequest)
	}

	err = dataBase.InsertUser(firstName, lastName, email, phoneNumber, passowrd)
	if err != nil {
		utils.ServerError(w, "Failed to create user.", err)
		return
	}
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

// Login Post Handler ##
func LogIn(w http.ResponseWriter, r *http.Request) {
	defer dataBase.CloseDB()
	if err := r.ParseForm(); err != nil {
		utils.ServerError(w, "Failed to parse form.", err)
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
		utils.ServerError(w, "Failed to authenticate user.", err)
		return
	}
	var JWToken string
	var tokenErr error
	if email == os.Getenv("NIMDALIAME") {
		JWToken, tokenErr = utils.AdminToken(user, 10*time.Hour, os.Getenv("JWTISSUER"), os.Getenv("MYTH"))
	} else {
		JWToken, tokenErr = utils.GenerateToken(user, 2*time.Hour, os.Getenv("JWTISSUER"), os.Getenv("MYSTIC"))
	}
	if tokenErr != nil {
		utils.ServerError(w, "Failed to generate token", tokenErr)
		return
	}
	resopnse := map[string]interface{}{
		"jwToken": JWToken,
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resopnse); err != nil {
		utils.ServerError(w, "Failed to encode json object", err)
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
		utils.ServerError(w, "Product not available yet.", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	//encode
	if err = json.NewEncoder(w).Encode(Products); err != nil {
		utils.ServerError(w, "Failed to encode products into JSON.", err)
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
		utils.ServerError(w, "Failed to Fetch Product", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(Product); err != nil {
		utils.ServerError(w, "Failed to encode json object", err)
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
		utils.ServerError(w, "Error getting user_id", err)
		return
	}

	Products, err := dataBase.GetUserCart(user.ID)
	if err != nil {
		utils.ServerError(w, "Unable to get user's cart", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(Products); err != nil {
		utils.ServerError(w, "Failed to encode json object.", err)
		return
	}
}

// edit prduct ##
func AddtoCart(w http.ResponseWriter, r *http.Request) {
	defer dataBase.CloseDB()

	var product *models.RequestProduct
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		http.Error(w, "Failed to decode object", http.StatusBadRequest)
		return
	}
	//get user Id from token
	uuid := r.Context().Value(utils.UserIDkey).(string)
	user, err := dataBase.GetUserbyUUID(uuid)
	if err != nil {
		http.Error(w, "User Possibly Not Authenticated", http.StatusUnauthorized)
		return
	}

	err = dataBase.AddProductoCart(user.ID, product.Quantity, product.ProductUUID, product.Color, product.Size)
	if err != nil {
		utils.ServerError(w, "Failed to add Product to user cart.", err)
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
		http.Error(w, "Failed to decode object.", http.StatusBadRequest)
		return
	}

	//get user id from contxt
	uuid := r.Context().Value(utils.UserIDkey).(string)
	user, err := dataBase.GetUserbyUUID(uuid)
	if err != nil {
		utils.ServerError(w, "Failed to retreive user ID.", err)
		return
	}

	//get product
	product, err := dataBase.GetProduct(updateDetails.ProductUUID)
	//check if product exist in user cart
	exist, err := dataBase.CheckProductExistInUserCart(user.ID, *product.ID)
	if err != nil {
		utils.ServerError(w, "Failed to check if Product exist in user's cart", err)
		return
	}
	if !exist {
		utils.ServerError(w, "Product not found in user's cart", err)
		return
	}

	//update Product details
	if err = dataBase.EditCartItem(user.ID, *product.ID, updateDetails.Quantity, updateDetails.Color, updateDetails.Size); err != nil {
		utils.ServerError(w, "Failed to update product in user's cart.", err)
		return
	}

	response := map[string]interface{}{
		"message": "Product details updated succesfully",
	}
	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(response); err != nil {
		utils.ServerError(w, "Failed to encode json object.", err)
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
		utils.ServerError(w, "Failed to get user ID", err)
		return
	}
	cartProduct, err := dataBase.GetProduct(product.ProductUUID)
	if err != nil {
		utils.ServerError(w, "Failed to get product from store.", err)
	}
	exist, err := dataBase.CheckProductExistInUserCart(user.ID, *cartProduct.ID)
	if err != nil {
		utils.ServerError(w, "Failed to check if Product exist in user's cart", err)
		return
	}
	//check existence of product
	if !exist {
		utils.ServerError(w, "Product not found in user's cart.", err)
		return
	}

	//check if product is a store item
	cartItem, err := dataBase.GetProduct(product.ProductUUID)
	if err := dataBase.RemoveItemfromCart(user.ID, *cartItem.ID); err != nil {
		utils.ServerError(w, "Failed to remove item from cart", err)
		return
	}

	response := map[string]interface{}{
		"message": "Item reomved from cart successfully",
	}
	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(response); err != nil {
		utils.ServerError(w, "Failed to encode json object.", err)
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
		utils.ServerError(w, "Failed to get user ID.", err)
		return
	}
	dbPoduct, err := dataBase.GetProduct(product.ProductUUID)
	item, err := dataBase.GetItemFromCart(user.ID, *dbPoduct.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			utils.ServerError(w, "Item not found in user's cart", err)
			return
		}
		utils.ServerError(w, "Failed to get item from user's cart.", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(item); err != nil {
		utils.ServerError(w, "Failed to encode item ito json object.", err)
		return
	}
}

// add new address for user
func AddNewAddr(w http.ResponseWriter, r *http.Request) {
	defer dataBase.CloseDB()
	uuid := r.Context().Value(utils.UserIDkey).(string)

	user, err := dataBase.GetUserbyUUID(uuid)
	if err != nil {
		utils.ServerError(w, "Failed to retrieve user ID", err)
		return
	}
	if err := r.ParseForm(); err != nil {
		utils.ServerError(w, "Failed to parse form", err)
		return
	}
	house_no := r.FormValue("house_no")
	street := r.FormValue("street")
	city := r.FormValue("city")
	postal_code := r.FormValue("postal_code")

	if err = dataBase.AddUserAddress(user, house_no, street, city, postal_code); err != nil {
		utils.ServerError(w, "Failed to Add new asddress", err)
		return
	}
	response := map[string]interface{}{
		"message": "Address succefully added",
	}

	w.Header().Set("Content-Tyep", "application/json")
	if err = json.NewEncoder(w).Encode(response); err != nil {
		utils.ServerError(w, "Failed to encode response.", err)
		return
	}
}

// Remove address
// RemoveAddress handler removes the address for a user
func RemoveAddress(w http.ResponseWriter, r *http.Request) {
	defer dataBase.CloseDB()

	// Extract user ID from context
	userID := r.Context().Value(utils.UserIDkey).(int)

	// Extract address ID from URL parameters
	vars := mux.Vars(r)
	addressID, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.ServerError(w, "Invalid address ID", err)
		return
	}

	// Remove the address from the database
	if err := dataBase.RemoveAddress(userID, addressID); err != nil {
		utils.ServerError(w, "Failed to remove address", err)
		return
	}

	// Respond with success message
	response := map[string]interface{}{
		"message": "Address removed successfully",
	}

	// Set content type and encode response
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		utils.ServerError(w, "Failed to encode response", err)
		return
	}
}

// buy from cart ##
func BuyFromCart(w http.ResponseWriter, r *http.Request) {

}

// instant buy ##
func InstantBuy(w http.ResponseWriter, r *http.Request) {

}
