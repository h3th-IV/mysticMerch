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
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

var (
	db, _    = database.InitDB()
	dataBase = database.DBModel{
		DB: db,
	}
)

// write api response in a go
func apiResponse(response map[string]interface{}, w http.ResponseWriter) {
	//set header
	w.Header().Set("Content-Type", "application/json")
	//decode json
	if err := json.NewEncoder(w).Encode(response); err != nil {
		utils.ServerError(w, "failed to encode object", err)
		return
	}
}

// home Handler display a list products
func Home(w http.ResponseWriter, r *http.Request) {
	//get some list of prduct to display on the home page
	products, err := dataBase.ViewProducts()
	if err != nil {
		utils.ReplaceLogger.Error("Failed to get product", zap.Error(err))
		utils.ServerError(w, "Failed to get products. Please try again later.", err)
		return
	}
	response := map[string]interface{}{
		"message": "Items retrived succesfully",
		"items":   products,
	}
	apiResponse(response, w)
	defer dataBase.CloseDB()
}

// admin stuff
func AddItemtoStore(w http.ResponseWriter, r *http.Request) {
	//decode new item
	var Product *models.Product
	if err := json.NewDecoder(r.Body).Decode(&Product); err != nil {
		utils.ReplaceLogger.Error("Failed to decode json", zap.Error(err))
		http.Error(w, "Failed to decode json", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	//add product to database
	_, err := dataBase.AddProduct(Product.ProductName, Product.Description, Product.Image, Product.Price)
	if err != nil {
		utils.ReplaceLogger.Error("Failed to add product", zap.Error(err))
		utils.ServerError(w, "Failed to add product to store", err)
		return
	}

	//write and send response
	response := map[string]interface{}{
		"message": "Operation was succesfull",
	}
	apiResponse(response, w)
	defer dataBase.CloseDB()
}

// admin stuff
func RemoveItemfromStore(w http.ResponseWriter, r *http.Request) {
	//decode item -- out of stock item
	var Product *models.RequestProduct
	if err := json.NewDecoder(r.Body).Decode(&Product); err != nil {
		utils.ReplaceLogger.Error("Failed to decode json", zap.Error(err))
		http.Error(w, "Failed to decode json object", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if err := dataBase.RemoveProductFromStore(Product.ProductUUID); err != nil {
		utils.ReplaceLogger.Error("Failed to remove item from store", zap.Error(err))
		utils.ServerError(w, "Failed to reomve itme from store", err)
		return
	}

	response := make(map[string]interface{})
	response["message"] = "Item Removed Succefully"
	apiResponse(response, w)
	defer dataBase.CloseDB()
}

// admin send mail ##
func AdminBroadcast(w http.ResponseWriter, r *http.Request) {
	//decode json object
	var notification *models.BroadcastNotification
	if err := json.NewDecoder(r.Body).Decode(&notification); err != nil {
		utils.ReplaceLogger.Error("Failed to decode json", zap.Error(err))
		http.Error(w, "Failed to decode json object", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	users, err := dataBase.GetAllUsers()
	if err != nil {
		utils.ReplaceLogger.Error("Failed to retrive users for broadcast message", zap.Error(err))
		http.Error(w, "Failed to retrive users for Brodcast message"+err.Error(), http.StatusInternalServerError)
		return
	}
	if err := admin.MarketingEmail(users, notification.Subject, notification.Body); err != nil {
		utils.ReplaceLogger.Error("Failed to send broadcast", zap.Error(err))
		http.Error(w, "Failed to send broadcast message"+err.Error(), http.StatusInternalServerError)
		return
	}
	response := map[string]interface{}{
		"message": "Broadcast email sent succesfully",
	}
	apiResponse(response, w)
	defer dataBase.CloseDB()
}

// send Transactional email to aparticular Customer
func Transactional(w http.ResponseWriter, r *http.Request) {
	var notification *models.TransactionNotification
	if err := json.NewDecoder(r.Body).Decode(&notification); err != nil {
		utils.ReplaceLogger.Error("Failed to decode json", zap.Error(err))
		http.Error(w, "Failed to decode json object", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if notification.ResponseUser.Email == "" || notification.Body == "" || notification.Subject == "" {
		http.Error(w, "User email, email body or email subject is empty", http.StatusInternalServerError)
		return
	}

	if err := admin.TransactionalEmail(&notification.ResponseUser, notification.Subject, notification.Body); err != nil {
		utils.ReplaceLogger.Error("Failed to send mail to usr", zap.Error(err))
		http.Error(w, "Failed to send email to user"+err.Error(), http.StatusInternalServerError)
		return
	}

	response := make(map[string]interface{})
	response["message"] = "email sent Successfully"

	apiResponse(response, w)
	defer dataBase.CloseDB()
}

// signUp post form Hadler ##
func SignUp(w http.ResponseWriter, r *http.Request) {
	//decode object
	var user *models.RequestUser
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		utils.ReplaceLogger.Error("Failed to decode json", zap.Error(err))
		http.Error(w, "Failed to decode json item", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	isDetails := utils.ValidateSignUpDetails([]models.ValidAta{
		{Value: user.FirstName, Validator: "firstname"},
		{Value: user.LastName, Validator: "lastname"},
		{Value: user.Email, Validator: "email"},
		{Value: user.Password, Validator: "password"},
	})

	//validate user input as w don't trust user input
	if !isDetails {
		http.Error(w, "Failed to Validate user details", http.StatusBadRequest)
		return
	}

	err := dataBase.InsertUser(user.FirstName, user.LastName, user.Email, user.PhoneNumber, user.Password)
	if err != nil {
		utils.ReplaceLogger.Error("Failed to create user", zap.Error(err))
		utils.ServerError(w, "Failed to create user.", err)
		return
	}
	response := map[string]interface{}{
		"message": "User account created succesffuly",
	}
	apiResponse(response, w)
	//http.Redirect(w, r, "/login", http.StatusSeeOther)
	defer dataBase.CloseDB()
}

// Login Post Handler ##
func LogIn(w http.ResponseWriter, r *http.Request) {
	//Load env var
	if err := utils.LoadEnv(); err != nil {
		utils.ReplaceLogger.Error("Failed to load env variables", zap.Error(err))
		http.Error(w, "Operation Failed", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	var Login *models.Login
	if err := json.NewDecoder(r.Body).Decode(&Login); err != nil {
		utils.ReplaceLogger.Error("Failed to decode json", zap.Error(err))
		http.Error(w, "failed to decode json object", http.StatusBadRequest)
		return
	}

	user, err := dataBase.AuthenticateUser(Login.Email)
	if err != nil {
		utils.ReplaceLogger.Error("Unable to retrieve user details", zap.Error(err))
		http.Error(w, "Unable to retrieve details", http.StatusUnauthorized)
		return
	}
	passErr := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(Login.Password))
	if passErr == bcrypt.ErrMismatchedHashAndPassword && passErr != nil {
		http.Error(w, "Password is incorect", http.StatusUnauthorized)
		return
	}
	var JWToken string
	var tokenErr error
	//create  admin token
	if Login.Email == os.Getenv("NIMDALIAME") {
		JWToken, tokenErr = utils.AdminToken(user, 10*time.Hour, os.Getenv("JWTISSUER"), os.Getenv("MYTH"))
	} else {
		//create user token
		JWToken, tokenErr = utils.GenerateToken(user, 2*time.Hour, os.Getenv("JWTISSUER"), os.Getenv("MYSTIC"))
	}
	if tokenErr != nil {
		utils.ServerError(w, "Failed to generate token", tokenErr)
		return
	}
	resopnse := map[string]interface{}{
		"message": "Login Succesfully",
		"jwToken": JWToken,
	}
	apiResponse(resopnse, w)
	defer dataBase.CloseDB()
}

// Serch product by query name
func SearchProduct(w http.ResponseWriter, r *http.Request) {
	//get name from search query
	query := r.URL.Query()
	ProductName := query.Get("product_name")

	Products, err := dataBase.GetProductByName(ProductName)
	if err != nil {
		utils.ReplaceLogger.Error("Product not found in user search", zap.Error(err))
		utils.ServerError(w, "Product not available yet.", err)
		return
	}
	response := map[string]interface{}{
		"message": "product found",
		"item":    Products,
	}
	apiResponse(response, w)
	defer dataBase.CloseDB()
}

// veiw product ##
func ViewProduct(w http.ResponseWriter, r *http.Request) {
	Var := mux.Vars(r)
	Product_id := Var["id"]

	//get by the uuid of product
	Product, err := dataBase.GetProduct(Product_id)
	if err != nil {
		utils.ReplaceLogger.Error("Failed to fetch product from DB", zap.Error(err))
		utils.ServerError(w, "Failed to Fetch Product", err)
		return
	}
	response := map[string]interface{}{
		"message": "Product details found",
		"item":    Product,
	}
	apiResponse(response, w)
	defer dataBase.CloseDB()
}

//Cart Operations

// view user cart ##
func GetUserCart(w http.ResponseWriter, r *http.Request) {
	//get user Id from token
	uuid := r.Context().Value(utils.UserIDkey).(string)
	user, err := dataBase.GetUserbyUUID(uuid)
	if err != nil {
		http.Error(w, "User Possibly Not Authenticated", http.StatusUnauthorized)
		return
	}

	//fecth cart
	Cart, err := dataBase.GetUserCart(user.ID)
	if err != nil {
		utils.ReplaceLogger.Error("Unable to fetch user's cart", zap.Error(err))
		utils.ServerError(w, "Unable to get user's cart", err)
		return
	}
	//write response
	response := map[string]interface{}{
		"message": "User cart returned succefully",
		"item":    Cart,
	}
	apiResponse(response, w)
	defer dataBase.CloseDB()
}

// edit prduct ##
func AddtoCart(w http.ResponseWriter, r *http.Request) {
	var product *models.RequestProduct
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		utils.ReplaceLogger.Error("Failed to decode json", zap.Error(err))
		http.Error(w, "failed to decode json object", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	//get user Id from token
	uuid := r.Context().Value(utils.UserIDkey).(string)
	user, err := dataBase.GetUserbyUUID(uuid)
	if err != nil {
		http.Error(w, "User Possibly Not Authenticated", http.StatusUnauthorized)
		return
	}

	err = dataBase.AddProductoCart(user.ID, product.Quantity, product.ProductUUID, product.Color, product.Size)
	if err != nil {
		utils.ReplaceLogger.Error("Failed to add product to cart", zap.Error(err))
		utils.ServerError(w, "Failed to add Product to user cart.", err)
		return
	}

	//write response
	response := make(map[string]interface{})
	response["message"] = "Product added to user cart succesfully"
	response["product"] = product

	apiResponse(response, w)
	defer dataBase.CloseDB()
}

// update cart details like add quantity, change color and size
func UpdateProductDetails(w http.ResponseWriter, r *http.Request) {
	//parse Update details
	var updateDetails *models.RequestProduct
	if err := json.NewDecoder(r.Body).Decode(&updateDetails); err != nil {
		http.Error(w, "failed to decode json object", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	//get user id from contxt
	uuid := r.Context().Value(utils.UserIDkey).(string)
	user, err := dataBase.GetUserbyUUID(uuid)
	if err != nil {
		utils.ServerError(w, "Failed to retreive user ID.", err)
		return
	}

	//get product
	product, err := dataBase.GetProduct(updateDetails.ProductUUID)
	if err != nil {
		utils.ReplaceLogger.Error("Failed to get product", zap.Error(err))
		utils.ServerError(w, "Failed to get product", err)
		return
	}
	//check if product exist in user cart
	exist, err := dataBase.CheckProductExistInUserCart(user.ID, product.ID)
	if err != nil {
		utils.ReplaceLogger.Error("Failed to check if product exist in user cart", zap.Error(err))
		utils.ServerError(w, "Failed to check if Product exist in user's cart", err)
		return
	}
	if !exist {
		utils.ReplaceLogger.Error("product does not exist in user cart", zap.Error(err))
		utils.ServerError(w, "Product not found in user's cart", err)
		return
	}

	//update Product details
	if err = dataBase.EditCartItem(user.ID, product.ID, updateDetails.Quantity, updateDetails.Color, updateDetails.Size); err != nil {
		utils.ReplaceLogger.Error("Failed to update  product details in user cart", zap.Error(err))
		utils.ServerError(w, "Failed to update product in user's cart.", err)
		return
	}

	response := map[string]interface{}{
		"message": "Product details updated succesfully",
	}
	apiResponse(response, w)
	defer dataBase.CloseDB()
}

// edit prduct ##
func RemovefromCart(w http.ResponseWriter, r *http.Request) {
	var product *models.RequestProduct
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		http.Error(w, "failed to decode json object", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	uuid := r.Context().Value(utils.UserIDkey).(string)
	user, err := dataBase.GetUserbyUUID(uuid)
	if err != nil {
		utils.ReplaceLogger.Error("Failed to get user id", zap.Error(err))
		utils.ServerError(w, "Failed to get user ID", err)
		return
	}
	cartProduct, err := dataBase.GetProduct(product.ProductUUID)
	if err != nil {
		utils.ReplaceLogger.Error("Failed to get product from store", zap.Error(err))
		utils.ServerError(w, "Failed to get product from store.", err)
		return
	}
	exist, err := dataBase.CheckProductExistInUserCart(user.ID, cartProduct.ID)
	if err != nil {
		utils.ReplaceLogger.Error("Failed to check if product exist in user cart", zap.Error(err))
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
	if err != nil {
		utils.ReplaceLogger.Error("Failed to get product", zap.Error(err))
		utils.ServerError(w, "failed to get product", err)
		return
	}
	if err := dataBase.RemoveItemfromCart(user.ID, cartItem.ID); err != nil {
		utils.ReplaceLogger.Error("Failed to remove item from cart", zap.Error(err))
		utils.ServerError(w, "failed to remove item from cart", err)
		return
	}

	response := map[string]interface{}{
		"message": "Item removed from cart successfully",
	}
	apiResponse(response, w)
	defer dataBase.CloseDB()
}

// GetItem from cart use list UserPorduct here ##
func GetItemFromCart(w http.ResponseWriter, r *http.Request) {
	var product *models.RequestProduct
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		http.Error(w, "failed to decode json object", http.StatusBadRequest)
		return
	}

	defer r.Body.Close()

	uuid := r.Context().Value(utils.UserIDkey).(string)
	user, err := dataBase.GetUserbyUUID(uuid)
	if err != nil {
		utils.ServerError(w, "Failed to get user ID.", err)
		return
	}
	dbPoduct, err := dataBase.GetProduct(product.ProductUUID)
	if err != nil {
		utils.ReplaceLogger.Error("Failed to  fetch product", zap.Error(err))
		utils.ServerError(w, "Failed to fetch product", err)
		return
	}
	item, err := dataBase.GetItemFromCart(user.ID, dbPoduct.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			utils.ReplaceLogger.Error("item not found in user cart", zap.Error(err))
			utils.ServerError(w, "Item not found in user's cart", err)
			return
		}
		utils.ServerError(w, "Failed to get item from user's cart.", err)
		return
	}
	//write response
	response := map[string]interface{}{
		"message": "Item retrived from user's cart",
		"item":    item,
	}
	apiResponse(response, w)
	defer dataBase.CloseDB()
}

// add new address for user
func AddNewAddr(w http.ResponseWriter, r *http.Request) {
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
		utils.ReplaceLogger.Error("Failed to add new address for user", zap.Error(err))
		utils.ServerError(w, "Failed to Add new asddress", err)
		return
	}
	response := map[string]interface{}{
		"message": "Address succesfully added",
	}
	apiResponse(response, w)
	defer dataBase.CloseDB()
}

// Remove address
// RemoveAddress handler removes the address for a user
func RemoveAddress(w http.ResponseWriter, r *http.Request) {
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
		utils.ReplaceLogger.Error("Failed to remove address", zap.Error(err))
		utils.ServerError(w, "Failed to remove address", err)
		return
	}

	// Respond with success message
	response := map[string]interface{}{
		"message": "Address removed successfully",
	}

	apiResponse(response, w)
	defer dataBase.CloseDB()
}

// buy from cart ##
func BuyFromCart(w http.ResponseWriter, r *http.Request) {

}

// instant buy ##
func InstantBuy(w http.ResponseWriter, r *http.Request) {

}
