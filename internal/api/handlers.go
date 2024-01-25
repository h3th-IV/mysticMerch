package api

import (
	"fmt"
	"net/http"

	"github.com/h3th-IV/mysticMerch/internal/database"
	"github.com/h3th-IV/mysticMerch/internal/models"
	"github.com/h3th-IV/mysticMerch/internal/utils"
)

var (
	db, err = database.InitDB()
	data    = database.DBModel{
		DB: db,
	}
)

// home Handler
func Home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Welcome to MysticeMerch")
}

// signUp post form Hadler ##
func SignUp(w http.ResponseWriter, r *http.Request) {
	defer data.CloseDB()
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
		{Value: firstName, Validator: "fName"},
		{Value: lastName, Validator: "lName"},
		{Value: email, Validator: "email"},
		{Value: passowrd, Validator: "password"},
	})
	if !isDetailsValid {
		http.Error(w, "Invalid User Input", http.StatusBadRequest)
	}

	err = data.InsertUser(firstName, lastName, email, phoneNumber, passowrd)
	if err != nil {
		utils.ServerError(w, err)
		return
	}
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

// Login Post Handler ##
func LogIn(w http.ResponseWriter, r *http.Request) {
	defer data.CloseDB()
	if err := r.ParseForm(); err != nil {
		utils.ServerError(w, err)
	}
	email := r.FormValue("email")
	password := r.FormValue("password")
	user, err := data.AuthenticateUser(email, password)
	if err != nil {
		utils.ServerError(w, fmt.Errorf("User"))
	}

}

// Serch product by query ##
func SearchProduct(w http.ResponseWriter, r *http.Request) {

}

// veiw product ##
func ViewProducts(w http.ResponseWriter, r *http.Request) {

}

//Cart Operations

// update product details like add quantity ##
func UpdateProductDetails(w http.ResponseWriter, r *http.Request) {

}

// view user cart ##
func UserCart(w http.ResponseWriter, r *http.Request) {

}

// edit prduct ##
func AddtoCart(w http.ResponseWriter, r *http.Request) {

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
