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
		http.Error(w, "failed to encode json object", http.StatusInternalServerError)
		return
	}
}

// home Handler display a list products
func Home(w http.ResponseWriter, r *http.Request) {
	//get some list of prduct to display on the home page
	products, err := dataBase.ViewProducts()
	if err != nil {
		utils.ReplaceLogger.Error("failed to get product", zap.Error(err))
		response := map[string]interface{}{
			"message": "failed to get products",
		}
		http.Error(w, "", http.StatusNotFound)
		apiResponse(response, w)
		return
	}
	response := map[string]interface{}{
		"message": "items retreived succesfully",
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
		response := map[string]interface{}{
			"message": "failed to add product to store",
		}
		http.Error(w, "", http.StatusInternalServerError)
		apiResponse(response, w)
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
		response := map[string]interface{}{
			"message": "failed to remove item from store",
		}
		http.Error(w, "", http.StatusInternalServerError)
		apiResponse(response, w)
		return
	}

	response := make(map[string]interface{})
	response["message"] = "item Removed Succesfully"
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
		utils.ReplaceLogger.Error("failed to create user account", zap.Error(err))
		response := map[string]interface{}{
			"message": "failed to create user account",
		}
		http.Error(w, "", http.StatusInternalServerError)
		apiResponse(response, w)
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
		utils.ReplaceLogger.Error("err generating token", zap.Error(tokenErr))
		http.Error(w, "error generating token", http.StatusInternalServerError)
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
		response := map[string]interface{}{
			"message": "product not available",
		}
		http.Error(w, "", http.StatusNotFound)
		apiResponse(response, w)
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
		response := map[string]interface{}{
			"message": "failed to fecth product from store",
		}
		http.Error(w, "", http.StatusInternalServerError)
		apiResponse(response, w)
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
		response := map[string]interface{}{
			"mesaage": "unable to fetch user's cart",
		}
		utils.ServerError(w, "unable to fetch user cart", err)
		apiResponse(response, w)
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
	if err != nil {
		utils.ReplaceLogger.Error("failed to retrieve product from store", zap.Error(err))
		http.Error(w, "failed to retreive product from store", http.StatusInternalServerError)
		return
	}
	if productExist != 1 {
		utils.ReplaceLogger.Error("product not found in store")
		response := map[string]interface{}{
			"response": "product not found in store",
		}
		http.Error(w, "", http.StatusNotFound)
		apiResponse(response, w)
	}

	err = dataBase.AddProductoCart(user.ID, product.Quantity, product.ProductUUID, product.Color, product.Size)
	if err != nil {
		utils.ReplaceLogger.Error("failed to add product to cart", zap.Error(err))
		response := map[string]interface{}{
			"response": "failed to add product to cart",
		}
		http.Error(w, "", http.StatusInternalServerError)
		apiResponse(response, w)
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
		utils.ReplaceLogger.Error("failed to retrieve user id", zap.Error(err))
		http.Error(w, "failed to retrieve user id", http.StatusInternalServerError)
		return
	}

	//get product
	product, err := dataBase.GetProduct(updateDetails.ProductUUID)
	if err != nil {
		utils.ReplaceLogger.Error("failed to get product", zap.Error(err))
		response := map[string]interface{}{
			"message": "failed to get product",
		}
		http.Error(w, "", http.StatusInternalServerError)
		apiResponse(response, w)
		return
	}

	//check if product exist in user cart
	exist, err := dataBase.CheckProductExistInUserCart(user.ID, product.ProductID)
	if err != nil {
		utils.ReplaceLogger.Error("failed to check if product exist in user cart", zap.Error(err))
		response := map[string]interface{}{
			"message": "failed to check if product exist in user cart",
		}
		http.Error(w, "", http.StatusInternalServerError)
		apiResponse(response, w)
		return
	}
	if !exist {
		utils.ReplaceLogger.Error("product does not exist in user cart", zap.Error(err))
		response := map[string]interface{}{
			"message": "product does not exist in user cart",
		}
		http.Error(w, "product not found in user cart", http.StatusNotFound)
		apiResponse(response, w)
		return
	}

	//update Product details
	if err = dataBase.EditCartItem(user.ID, product.ID, updateDetails.Quantity, updateDetails.Color, updateDetails.Size); err != nil {
		utils.ReplaceLogger.Error("failed to update product details", zap.Error(err))
		response := map[string]interface{}{
			"message": "failed to update product details",
		}
		http.Error(w, "", http.StatusInternalServerError)
		apiResponse(response, w)
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
		http.Error(w, "failed to get user id", http.StatusInternalServerError)

		return
	}

	//check if product in store...don't hate me just sayin
	instore, err := dataBase.CheckProductExist(product.ProductUUID)
	if err != nil {
		utils.ReplaceLogger.Error("failed to check if product not in user store", zap.Error(err))
		http.Error(w, "failed to check if product in store", http.StatusInternalServerError)
		return
	}
	if instore != 1 {
		utils.ReplaceLogger.Error("product not found in store", zap.Error(err))
		response := map[string]interface{}{
			"message": "product not found in store",
		}
		http.Error(w, "product not found in store", http.StatusNotFound)
		apiResponse(response, w)
		return
	}

	//check if product in user cart
	exist, err := dataBase.CheckProductExistInUserCart(user.ID, product.ProductUUID)
	if err != nil {
		utils.ReplaceLogger.Error("failed to check if product exist in user cart", zap.Error(err))
		response := map[string]interface{}{
			"message": "failed to check if product not in user store",
		}
		http.Error(w, "", http.StatusInternalServerError)
		apiResponse(response, w)
		return
	}
	//check existence of product
	if !exist {
		response := map[string]interface{}{
			"message": "product not found in user's cart",
		}
		http.Error(w, "product not found in user's cart", http.StatusNotFound)
		apiResponse(response, w)
		return
	}

	if err := dataBase.RemoveItemfromCart(user.ID, product.ProductUUID); err != nil {
		utils.ReplaceLogger.Error("failed to remove item from user's cart", zap.Error(err))
		response := map[string]interface{}{
			"message": "failed to remove item from user's cart",
		}
		http.Error(w, "", http.StatusInternalServerError)
		apiResponse(response, w)
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
		http.Error(w, "failed to get user id", http.StatusInternalServerError)
		return
	}
	dbPoduct, err := dataBase.GetProduct(product.ProductUUID)
	if err != nil {
		utils.ReplaceLogger.Error("failed to  fetch product", zap.Error(err))
		response := map[string]interface{}{
			"message": "failed to fetch product",
		}
		http.Error(w, "", http.StatusInternalServerError)
		apiResponse(response, w)
		return
	}
	item, err := dataBase.GetItemFromCart(user.ID, dbPoduct.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			utils.ReplaceLogger.Error("item not found in user cart", zap.Error(err))
			response := map[string]interface{}{
				"message": "item not found in user cart",
			}
			http.Error(w, "", http.StatusNotFound)
			apiResponse(response, w)
			return
		}
		http.Error(w, "", http.StatusInternalServerError)
		response := map[string]interface{}{
			"message": "failed to get item from user's cart",
		}
		apiResponse(response, w)
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
		http.Error(w, "failed to retrieve user id", http.StatusInternalServerError)
		return
	}
	if err := r.ParseForm(); err != nil {
		http.Error(w, "failed to parse form", http.StatusInternalServerError)
		return
	}
	house_no := r.FormValue("house_no")
	street := r.FormValue("street")
	city := r.FormValue("city")
	postal_code := r.FormValue("postal_code")

	if err = dataBase.AddUserAddress(user, house_no, street, city, postal_code); err != nil {
		utils.ReplaceLogger.Error("failed to add new address for user", zap.Error(err))
		response := map[string]interface{}{
			"message": "failed to add new address",
		}
		http.Error(w, "", http.StatusInternalServerError)
		apiResponse(response, w)
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
		http.Error(w, "invalid address id", http.StatusBadRequest)
		return
	}

	// Remove the address from the database
	if err := dataBase.RemoveAddress(userID, addressID); err != nil {
		utils.ReplaceLogger.Error("failed to remove address", zap.Error(err))
		response := map[string]interface{}{
			"message": "failed to remove address",
		}
		http.Error(w, "", http.StatusInternalServerError)
		apiResponse(response, w)
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
