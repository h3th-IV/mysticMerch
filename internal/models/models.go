package models

import (
	"time"
)

// struct for the marketplace

// Usser model.
type User struct {
	ID             *int           `json:"id"`
	FirstName      *string        `json:"firstName" validate:"required,min=2,max=50"`
	LastName       *string        `json:"lastName" validate:"required,min=2,max=50"`
	Password       *string        `json:"password"`
	Email          *string        `json:"email" validate:"required,email"`
	PhoneNumber    *string        `json:"phoneNumber" validate:"required"`
	CreatedAt      time.Time      `json:"createdAt"`
	UpdatedAt      time.Ticker    `json:"updatedAt"`
	UserID         string         `json:"userId"`
	UserCart       []UserProducts `json:"userCart"`
	AddressDetails []Address      `json:"addressDetails"`
	OrderStatus    []Order        `json:"orderStatus"`
}

// Products available in store.
type Product struct {
	ProductID   *int    `json:"productId"`
	ProductName *string `json:"productName"`
	Price       *uint64 `json:"price"`
	Rating      *uint8  `json:"rating"`
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

// user's address details.
type Address struct {
	AddressID  *int    `json:"addressId"`
	HouseNo    *string `json:"houseNo"`
	Street     *string `json:"street"`
	City       *string `json:"city"`
	PostalCode *string `json:"postalCode"`
}

// Oorder model
type Order struct {
	OrderID       *int           `json:"orderId"`
	OrderCart     []UserProducts `json:"orderCart"`
	OrderedAt     time.Time      `json:"orderAt"`
	Price         *int           `json:"orderPrice"`
	Discount      *int           `json:"discount"`
	PaymentMethod Payment        `json:"paymentType"`
}

// payment method for an order, indicating whether electronic payment or cash was used.
type Payment struct {
	EletronicPayment bool `json:"electronicPayment"`
	Cash             bool `json:"cash"`
}
