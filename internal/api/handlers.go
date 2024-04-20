package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
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
		utils.ReplaceLogger.Error("failed to get product", zap.Error(err))
		utils.ServerError(w, "failed to get products. Please try again later.", err)
		return
	}
	response := map[string]interface{}{
		"message": "items retrived succesfully",
		"items":   products,
	}
	apiResponse(response, w)
}

// admin stuff
func AddItemtoStore(w http.ResponseWriter, r *http.Request) {
	//get Admin id
	uuid := r.Context().Value(utils.UserIDkey).(string)
	user, err := dataBase.GetUserbyUUID(uuid)
	if err != nil {
		http.Error(w, "user not authenticated", http.StatusNetworkAuthenticationRequired)
		return
	}
	if user.ID != 1 {
		http.Error(w, "user not authorised", http.StatusUnauthorized)
		return
	}
	//decode new item
	var Product *models.NewProduct
	if err := json.NewDecoder(r.Body).Decode(&Product); err != nil {
		utils.ReplaceLogger.Error("failed to decode json", zap.Error(err))
		http.Error(w, "failed to decode json", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	//add product to database
	_, err = dataBase.AddProduct(user.ID, Product.ProductName, Product.Description, Product.Image, Product.Price)
	if err != nil {
		utils.ReplaceLogger.Error("failed to add product", zap.Error(err))
		utils.ServerError(w, "failed to add product to store", err)
		return
	}

	//write and send response
	response := map[string]interface{}{
		"message": "operation was succesfull",
	}
	apiResponse(response, w)
}

// admin stuff
func RemoveItemfromStore(w http.ResponseWriter, r *http.Request) {
	//get usr id to cofirm if admin
	uuid := r.Context().Value(utils.UserIDkey).(string)
	user, err := dataBase.GetUserbyUUID(uuid)
	if err != nil {
		http.Error(w, "user not authenticated", http.StatusNetworkAuthenticationRequired)
		return
	}

	if user.ID != 1 {
		http.Error(w, "user not authorized", http.StatusUnauthorized)
		return
	}
	//decode item -- out of stock item
	var Product *models.RemoveProduct
	if err := json.NewDecoder(r.Body).Decode(&Product); err != nil {
		utils.ReplaceLogger.Error("failed to decode json", zap.Error(err))
		http.Error(w, "failed to decode json object", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if err := dataBase.RemoveProductFromStore(Product.ProductUUID); err != nil {
		utils.ReplaceLogger.Error("failed to remove item from store", zap.Error(err))
		utils.ServerError(w, "failed to reomve itme from store", err)
		return
	}

	response := make(map[string]interface{})
	response["message"] = "item Removed Succefully"
	apiResponse(response, w)
}

// admin send mail ##
func AdminBroadcast(w http.ResponseWriter, r *http.Request) {
	//get admin id
	uuid := r.Context().Value(utils.UserIDkey).(string)
	user, err := dataBase.GetUserbyUUID(uuid)
	if err != nil {
		http.Error(w, "user not authenticated", http.StatusNetworkAuthenticationRequired)
		return
	}
	if user.ID != 1 {
		http.Error(w, "user not authorised", http.StatusUnauthorized)
		return
	}
	//decode json object
	var notification *models.BroadcastNotification
	if err := json.NewDecoder(r.Body).Decode(&notification); err != nil {
		utils.ReplaceLogger.Error("failed to decode json", zap.Error(err))
		http.Error(w, "failed to decode json object", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	users, err := dataBase.GetAllUsers()
	if err != nil {
		utils.ReplaceLogger.Error("failed to retrive users for broadcast message", zap.Error(err))
		http.Error(w, "failed to retrive users for brodcast message"+err.Error(), http.StatusInternalServerError)
		return
	}
	if err := admin.MarketingEmail(users, notification.Subject, notification.Body); err != nil {
		utils.ReplaceLogger.Error("failed to send broadcast", zap.Error(err))
		http.Error(w, "failed to send broadcast message"+err.Error(), http.StatusInternalServerError)
		return
	}
	response := map[string]interface{}{
		"message": "broadcast email sent succesfully",
	}
	apiResponse(response, w)
}

// send Transactional email to aparticular Customer
func Transactional(w http.ResponseWriter, r *http.Request) {
	//auth admin
	uuid := r.Context().Value(utils.UserIDkey).(string)
	user, err := dataBase.GetUserbyUUID(uuid)
	if err != nil {
		http.Error(w, "user not authenticated", http.StatusNetworkAuthenticationRequired)
		return
	}
	if user.ID != 1 {
		http.Error(w, "user not authorized", http.StatusUnauthorized)
		return
	}
	var notification *models.TransactionNotification
	if err := json.NewDecoder(r.Body).Decode(&notification); err != nil {
		utils.ReplaceLogger.Error("failed to decode json", zap.Error(err))
		http.Error(w, "failed to decode json object", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if notification.ResponseUser.Email == "" || notification.Body == "" || notification.Subject == "" {
		http.Error(w, "user email, email body or email subject is empty", http.StatusInternalServerError)
		return
	}

	if err := admin.TransactionalEmail(&notification.ResponseUser, notification.Subject, notification.Body); err != nil {
		utils.ReplaceLogger.Error("failed to send mail to usr", zap.Error(err))
		http.Error(w, "failed to send email to user"+err.Error(), http.StatusInternalServerError)
		return
	}

	response := make(map[string]interface{})
	response["message"] = "email sent successfully"

	apiResponse(response, w)
}

// signUp post form Hadler ##
func SignUp(w http.ResponseWriter, r *http.Request) {
	//decode object
	var user *models.RequestUser
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		utils.ReplaceLogger.Error("failed to decode json", zap.Error(err))
		http.Error(w, "failed to decode json item", http.StatusBadRequest)
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
		http.Error(w, "failed to validate user details", http.StatusBadRequest)
		return
	}

	err := dataBase.InsertUser(user.FirstName, user.LastName, user.Email, user.PhoneNumber, user.Password)
	if err != nil {
		utils.ReplaceLogger.Error("failed to create user", zap.Error(err))
		utils.ServerError(w, "failed to create user.", err)
		return
	}
	response := map[string]interface{}{
		"message": "user account created succesffuly",
	}
	apiResponse(response, w)
	//http.Redirect(w, r, "/login", http.StatusSeeOther)
}

// Login Post Handler ##
func LogIn(w http.ResponseWriter, r *http.Request) {
	//Load env var
	if err := utils.LoadEnv(); err != nil {
		utils.ReplaceLogger.Error("failed to load env variables", zap.Error(err))
		http.Error(w, "operation Failed", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	var Login *models.Login
	if err := json.NewDecoder(r.Body).Decode(&Login); err != nil {
		utils.ReplaceLogger.Error("failed to decode json", zap.Error(err))
		http.Error(w, "failed to decode json object", http.StatusBadRequest)
		return
	}

	user, err := dataBase.AuthenticateUser(Login.Email)
	if err != nil {
		utils.ReplaceLogger.Error("unable to retrieve user details", zap.Error(err))
		http.Error(w, "unable to retrieve details", http.StatusUnauthorized)
		return
	}
	passErr := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(Login.Password))
	if passErr == bcrypt.ErrMismatchedHashAndPassword && passErr != nil {
		http.Error(w, "password is incorrect", http.StatusUnauthorized)
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
		utils.ServerError(w, "failed to generate token", tokenErr)
		return
	}
	resopnse := map[string]interface{}{
		"message": "login Succesfully",
		"jwToken": JWToken,
	}
	apiResponse(resopnse, w)
}

// Serch product by query name
func SearchProduct(w http.ResponseWriter, r *http.Request) {
	//get name from search query
	query := r.URL.Query()
	ProductName := query.Get("product_name")

	Products, err := dataBase.GetProductByName(ProductName)
	if err != nil {
		utils.ReplaceLogger.Error("product not found in user search", zap.Error(err))
		utils.ServerError(w, "product not available yet.", err)
		return
	}
	response := map[string]interface{}{
		"message": "product found",
		"item":    Products,
	}
	apiResponse(response, w)
}

// veiw product ##
func ViewProduct(w http.ResponseWriter, r *http.Request) {
	Var := mux.Vars(r)
	Product_id := Var["id"]

	//get by the uuid of product
	Product, err := dataBase.GetProduct(Product_id)
	if err != nil {
		utils.ReplaceLogger.Error("failed to fetch product from DB", zap.Error(err))
		utils.ServerError(w, "failed to fetch product", err)
		return
	}
	response := map[string]interface{}{
		"message": "product details found",
		"item":    Product,
	}
	apiResponse(response, w)
}

//Cart Operations

// view user cart ##
func GetUserCart(w http.ResponseWriter, r *http.Request) {
	//get user Id from token
	uuid := r.Context().Value(utils.UserIDkey).(string)
	user, err := dataBase.GetUserbyUUID(uuid)
	if err != nil {
		http.Error(w, "user possibly not authenticated", http.StatusNetworkAuthenticationRequired)
		return
	}

	//fecth cart
	Cart, err := dataBase.GetUserCart(user.ID)
	if err != nil {
		utils.ReplaceLogger.Error("unable to fetch user's cart", zap.Error(err))
		utils.ServerError(w, "unable to get user's cart", err)
		return
	}
	//write response
	response := map[string]interface{}{
		"message": "user cart returned succefully",
		"item":    Cart,
	}
	apiResponse(response, w)
}

// edit prduct ##
func AddtoCart(w http.ResponseWriter, r *http.Request) {
	//get user Id from token
	uuid := r.Context().Value(utils.UserIDkey).(string)
	user, err := dataBase.GetUserbyUUID(uuid)
	if err != nil {
		http.Error(w, "user possibly not authenticated", http.StatusUnauthorized)
		return
	}

	var product *models.RequestProduct
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		utils.ReplaceLogger.Error("failed to decode json", zap.Error(err))
		http.Error(w, "failed to decode json object", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	productExist, err := dataBase.CheckProductExist(product.ProductUUID)
	fmt.Println(productExist)
	if productExist != 1 {
		utils.ReplaceLogger.Error("failed to retrive product, product may not exist in store")
		utils.ServerError(w, "failed to retrive product, product may not exist in store", err)
	}
	if err != nil {
		utils.ReplaceLogger.Error("failed to retrieve product from store, product might not exist", zap.Error(err))
		utils.ServerError(w, "failed to retrieve product from store, product might not exists", err)
		return
	}

	err = dataBase.AddProductoCart(user.ID, product.Quantity, product.ProductUUID, product.Color, product.Size)
	if err != nil {
		utils.ReplaceLogger.Error("failed to add product to cart", zap.Error(err))
		utils.ServerError(w, "failed to add product to user cart.", err)
		return
	}

	//write response
	response := make(map[string]interface{})
	response["message"] = "product added to user cart succesfully"
	response["product"] = product

	apiResponse(response, w)
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
		utils.ServerError(w, "failed to retreive user id.", err)
		return
	}

	//get product
	product, err := dataBase.GetProduct(updateDetails.ProductUUID)
	if err != nil {
		utils.ReplaceLogger.Error("failed to get product", zap.Error(err))
		utils.ServerError(w, "failed to get product", err)
		return
	}
	//check if product exist in user cart
	exist, err := dataBase.CheckProductExistInUserCart(user.ID, product.ProductID)
	if err != nil {
		utils.ReplaceLogger.Error("failed to check if product exist in user cart", zap.Error(err))
		utils.ServerError(w, "failed to check if Product exist in user's cart", err)
		return
	}
	if !exist {
		utils.ReplaceLogger.Error("product does not exist in user cart", zap.Error(err))
		utils.ServerError(w, "product not found in user's cart", err)
		return
	}

	//update Product details
	if err = dataBase.EditCartItem(user.ID, product.ID, updateDetails.Quantity, updateDetails.Color, updateDetails.Size); err != nil {
		utils.ReplaceLogger.Error("failed to update  product details in user cart", zap.Error(err))
		utils.ServerError(w, "failed to update product in user's cart.", err)
		return
	}

	response := map[string]interface{}{
		"message": "product details updated succesfully",
	}
	apiResponse(response, w)
}

// remove product from user cart
func RemovefromCart(w http.ResponseWriter, r *http.Request) {
	var product *models.RemoveProduct
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		http.Error(w, "failed to decode json object", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	uuid := r.Context().Value(utils.UserIDkey).(string)
	user, err := dataBase.GetUserbyUUID(uuid)
	if err != nil {
		utils.ReplaceLogger.Error("failed to get user id", zap.Error(err))
		utils.ServerError(w, "failed to get user id", err)
		return
	}

	//check if product in store...don't hate me just sayin
	instore, err := dataBase.CheckProductExist(product.ProductUUID)
	if err != nil {
		utils.ReplaceLogger.Error("failed to check if product not in user store", zap.Error(err))
		utils.ServerError(w, "failed to check if product not in user store", err)
		return
	}
	if instore != 1 {
		utils.ReplaceLogger.Error("product not found in store", zap.Error(err))
		utils.ServerError(w, "product not found in store", err)
		return
	}
	//check if product in user cart
	exist, err := dataBase.CheckProductExistInUserCart(user.ID, product.ProductUUID)
	if err != nil {
		utils.ReplaceLogger.Error("failed to check if product exist in user cart", zap.Error(err))
		utils.ServerError(w, "failed to check if Product exist in user's cart", err)
		return
	}
	//check existence of product
	if !exist {
		utils.ServerError(w, "product not found in user's cart.", err)
		return
	}

	if err := dataBase.RemoveItemfromCart(user.ID, product.ProductUUID); err != nil {
		utils.ReplaceLogger.Error("failed to remove item from cart", zap.Error(err))
		utils.ServerError(w, "failed to remove item from cart", err)
		return
	}

	response := map[string]interface{}{
		"message": "item removed from cart successfully",
	}
	apiResponse(response, w)
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
		utils.ServerError(w, "failed to get user id.", err)
		return
	}
	dbPoduct, err := dataBase.GetProduct(product.ProductUUID)
	if err != nil {
		utils.ReplaceLogger.Error("failed to  fetch product", zap.Error(err))
		utils.ServerError(w, "failed to fetch product", err)
		return
	}
	item, err := dataBase.GetItemFromCart(user.ID, dbPoduct.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			utils.ReplaceLogger.Error("item not found in user cart", zap.Error(err))
			utils.ServerError(w, "item not found in user's cart", err)
			return
		}
		utils.ServerError(w, "failed to get item from user's cart.", err)
		return
	}
	//write response
	response := map[string]interface{}{
		"message": "item retrived from user's cart",
		"item":    item,
	}
	apiResponse(response, w)
}

// add new address for user
func AddNewAddr(w http.ResponseWriter, r *http.Request) {
	uuid := r.Context().Value(utils.UserIDkey).(string)

	user, err := dataBase.GetUserbyUUID(uuid)
	if err != nil {
		utils.ServerError(w, "failed to retrieve user id", err)
		return
	}
	if err := r.ParseForm(); err != nil {
		utils.ServerError(w, "failed to parse form", err)
		return
	}
	house_no := r.FormValue("house_no")
	street := r.FormValue("street")
	city := r.FormValue("city")
	postal_code := r.FormValue("postal_code")

	if err = dataBase.AddUserAddress(user, house_no, street, city, postal_code); err != nil {
		utils.ReplaceLogger.Error("failed to add new address for user", zap.Error(err))
		utils.ServerError(w, "failed to add new address", err)
		return
	}
	response := map[string]interface{}{
		"message": "address succesfully added",
	}
	apiResponse(response, w)
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
		utils.ServerError(w, "invalid address id", err)
		return
	}

	// Remove the address from the database
	if err := dataBase.RemoveAddress(userID, addressID); err != nil {
		utils.ReplaceLogger.Error("failed to remove address", zap.Error(err))
		utils.ServerError(w, "failed to remove address", err)
		return
	}

	// Respond with success message
	response := map[string]interface{}{
		"message": "address removed successfully",
	}

	apiResponse(response, w)
}

// buy from cart ##
func BuyFromCart(w http.ResponseWriter, r *http.Request) {

}

// instant buy ##
func InstantBuy(w http.ResponseWriter, r *http.Request) {

}

// func LogOut(w http.ResponseWriter, r *http.Request) {
// 	defer dataBase.CloseDB()

// 	r.Context().Value(utils.UserIDkey)
// 	c := http.Cookie{}
// }
