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
	ID          *int        `json:"id"`     //db auto increment
	UserID      string      `json:"userId"` //uuid
	FirstName   *string     `json:"firstName" validate:"required,min=2,max=50"`
	LastName    *string     `json:"lastName" validate:"required,min=2,max=50"`
	Email       *string     `json:"email" validate:"required,email"`
	PhoneNumber *string     `json:"phoneNumber" validate:"required"`
	Password    string      `json:"password"`
	CreatedAt   time.Time   `json:"createdAt"`
	UpdatedAt   time.Ticker `json:"updatedAt"`
}

// Products available in store.
type Product struct {
	ID          *int    `json:"id"`        //auto increment
	ProductID   *string `json:"productId"` //for non db ops uuid generated
	ProductName *string `json:"productName"`
	Description *string `json:"description"`
	Image       *string `json:"image"`
	Price       *uint64 `json:"price"`
	Rating      uint8   `json:"rating"`
}

// simplified product for API response
type ResponseProduct struct {
	ProductName *string `json:"productName"`
	Description *string `json:"description"`
	Price       *string `json:"price"`
	Rating      *string `json:"rating"`
	Image       *string `json:"image"`
}

// Produts associated with the user(like ordered product)
type UserProducts struct {
	ProductID   *int    `json:"productId"`
	ProductName *string `json:"productName"`
	Price       int     `json:"price"`
	Rating      *uint   `json:"rating"`
	Image       *string `json:"image"`
	Quantity    *int    `json:"quantity"`
	Color       *string `json:"color,omitempty"`
	Size        *string `json:"size,omitempty"`
}

// simplified cartProducts for API response
type ResponseCartProducts struct {
	ProductName *string `json:"productName"`
	Price       *int    `json:"price"`
	Rating      *uint   `json:"rating"`
	Image       *string `json:"image"`
	Quantity    *int    `json:"quantity"`
	Color       *string `json:"color,omitempty"`
	Size        *string `json:"sze,omitempty"`
}

// Oorder model
type Order struct {
	OrderID       *int      `json:"orderId"`
	OrderedAt     time.Time `json:"orderAt"`
	Price         *int      `json:"orderPrice"`
	Discount      *int      `json:"discount"`
	PaymentMethod Payment   `json:"paymentType"`
}

// user's address details.
type Address struct {
	AddressID   *int    `json:"addressId"`
	HouseNo     *string `json:"houseNo"`
	Street      *string `json:"street"`
	City        *string `json:"city"`
	PostalCode  *string `json:"postalCode"`
	UserPhoneNo *string `json:"phoneNumber"`
}

// payment method for an order, indicating whether electronic payment or cash was used.
type Payment struct {
	EletronicPayment bool `json:"electronicPayment"`
	Cash             bool `json:"cash"`
}

type ValidAta struct {
	Value     string
	Validator string
}
