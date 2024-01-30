package models

import (
	"time"
)

// struct for the marketplace
/*
	User
	Product
	UserProducts
	AdressTable
	Ordertable
	Payment
*/
// Usser model.
type User struct {
	ID          *int        `json:"id"`      //db auto increment
	UserID      string      `json:"user_id"` //uuid
	FirstName   *string     `json:"first_name" validate:"required,min=2,max=50"`
	LastName    *string     `json:"last_name" validate:"required,min=2,max=50"`
	Email       *string     `json:"email" validate:"required,email"`
	PhoneNumber *string     `json:"phone_number" validate:"required"`
	Password    string      `json:"password"`
	CreatedAt   time.Time   `json:"createdAt"`
	UpdatedAt   time.Ticker `json:"updatedAt"`
}

// // init new user type
// func NewUser(firstName, lastName, email, phoneNumber, password string) (*User, error) {
// 	uuid, err := utils.GenerateUUID("user")
// 	if err != nil {
// 		return nil, err
// 	}
// 	cryptedPassword, err := utils.EncryptPass([]byte(password))
// 	if err != nil {
// 		return nil, err
// 	}
// 	crypted := string(cryptedPassword)

// 	return &User{
// 		UserID:      uuid,
// 		FirstName:   &firstName,
// 		LastName:    &lastName,
// 		Email:       &email,
// 		PhoneNumber: &phoneNumber,
// 		Password:    crypted,
// 	}, nil
// }

// Products available in store.
type Product struct {
	ID          *int    `json:"id"`         //auto increment
	ProductID   *string `json:"product_id"` //for non db ops uuid generated
	ProductName *string `json:"product_name"`
	Description *string `json:"description"`
	Image       *string `json:"image"`
	Price       *uint64 `json:"price"`
	Rating      uint8   `json:"rating"`
}

// func NewProduct(name, description, image string, price uint64) (*Product, error) {
// 	uuid, err := utils.GenerateUUID("product")
// 	return &Product{
// 		ProductID:   &uuid,
// 		ProductName: &name,
// 		Description: &description,
// 		Image:       &image,
// 		Price:       &price,
// 		Rating:      uint8(0),
// 	}, err
// }

// simplified product for API response
type ResponseProduct struct {
	ProductName *string `json:"product_name"`
	Description *string `json:"description"`
	Price       *string `json:"price"`
	Rating      *string `json:"rating"`
	Image       *string `json:"image"`
}

// Produts associated with the user(like ordered product)
type UserProducts struct {
	ProductID   *int    `json:"product_id"`
	ProductName *string `json:"product_name"`
	Price       int     `json:"price"`
	Rating      *uint   `json:"rating"`
	Image       *string `json:"image"`
	Quantity    *int    `json:"quantity"`
	Color       *string `json:"color,omitempty"`
	Size        *string `json:"size,omitempty"`
}

// simplified cartProducts for API response
type ResponseCartProducts struct {
	ProductName *string `json:"product_name"`
	Price       *int    `json:"price"`
	Rating      *uint   `json:"rating"`
	Image       *string `json:"image"`
	Quantity    *int    `json:"quantity"`
	Color       *string `json:"color,omitempty"`
	Size        *string `json:"sze,omitempty"`
}

// Oorder model
type Order struct {
	OrderID       *int      `json:"order_id"`
	OrderedAt     time.Time `json:"order_at"`
	Price         *int      `json:"order_price"`
	Discount      *int      `json:"discount"`
	PaymentMethod Payment   `json:"payment_type"`
}

// user's address details.
type Address struct {
	AddressID   *int    `json:"address_id"`
	HouseNo     *string `json:"house_no"`
	Street      *string `json:"street"`
	City        *string `json:"city"`
	PostalCode  *string `json:"postal_code"`
	UserPhoneNo *string `json:"phone_number"`
}

// payment method for an order, indicating whether electronic payment or cash was used.
type Payment struct {
	EletronicPayment bool `json:"electronic_payment"`
	Cash             bool `json:"cash"`
}

type ValidAta struct {
	Value     string
	Validator string
}
