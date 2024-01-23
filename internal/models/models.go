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
	ID           *int        `json:"id"`
	FirstName    *string     `json:"firstName" validate:"required,min=2,max=50"`
	LastName     *string     `json:"lastName" validate:"required,min=2,max=50"`
	PasswordHash []byte      `json:"password"`
	Email        *string     `json:"email" validate:"required,email"`
	PhoneNumber  *string     `json:"phoneNumber" validate:"required"`
	CreatedAt    time.Time   `json:"createdAt"`
	UpdatedAt    time.Ticker `json:"updatedAt"`
	UserID       string      `json:"userId"`
}

// Products available in store.
type Product struct {
	ProductID   *int64  `json:"productId"`
	ProductName *string `json:"productName"`
	Description *string `json:"description"`
	Price       *uint64 `json:"price"`
	Rating      *uint8  `json:"rating"`
	Image       *string `json:"image"`
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
	AddressID  *int    `json:"addressId"`
	HouseNo    *string `json:"houseNo"`
	Street     *string `json:"street"`
	City       *string `json:"city"`
	PostalCode *string `json:"postalCode"`
}

// payment method for an order, indicating whether electronic payment or cash was used.
type Payment struct {
	EletronicPayment bool `json:"electronicPayment"`
	Cash             bool `json:"cash"`
}

type ValidAta struct {
	Value *string
	Valid *string
}
