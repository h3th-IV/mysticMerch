package models

import "time"

// Usser model.
type User struct {
	ID             *int
	FirstName      *string
	LastName       *string
	PasswordHash   *string
	Email          *string
	PhoneNumber    *string
	CreatedAt      time.Time
	UpdatedAt      time.Ticker
	UserID         string
	UserCart       []UserProducts
	AddressDetails []Address
	OrderStatus    []Order
}

// Products available in store.
type Product struct {
	ProductID   *int
	ProductName *string
	Price       *uint64
	Rating      *uint8
	Image       *string
}

// Produts associated with the user(like ordered product)
type UserProducts struct {
	ProductID   *int
	ProductName *string
	Price       int
	Rating      *uint
	Image       *string
}

// user's address details.
type Address struct {
	AddressID  *int
	HouseNo    *string
	Street     *string
	City       *string
	PostalCode *string
}

// Oorder model
type Order struct {
	OrderID       *int
	OrderCart     []UserProducts
	OrderedAt     time.Time
	Price         *int
	Discount      *int
	PaymentMethod Payment
}

// payment method for an order, indicating whether electronic payment or cash was used.
type Payment struct {
	EletronicPayment bool
	Cash             bool
}
