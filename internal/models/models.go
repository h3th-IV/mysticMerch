package models

import "time"

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

type Product struct {
	ProductID   *int
	ProductName *string
	Price       *uint64
	Rating      *uint8
	Image       *string //wil hold url
}
type UserProducts struct {
	ProductID   *int
	productName *string
	Price       int
	Rating      *uint
	Image       *string
}

type Address struct {
	AddressID  *int
	HouseNo    *string
	Street     *string
	City       *string
	PostalCode *string
}

type Order struct {
	OrderID       *int
	OrderCart     []UserProducts
	OrderedAt     time.Time
	Price         *int
	Discount      *int
	PaymentMethod Payment
}

type Payment struct {
	EPayment bool
	Cash     bool
}
