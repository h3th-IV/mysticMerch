package database

import (
	"database/sql"

	"github.com/h3th-IV/mysticMerch/internal/models"
	"github.com/h3th-IV/mysticMerch/internal/utils"
)

type ProductModel struct {
	DB *sql.DB
}

func NewProduct(name, description, image string, price uint64) *models.Product {
	uuid, _ := utils.GenerateUUID("product")
	return &models.Product{
		ProductID:   &uuid,
		ProductName: &name,
		Description: &description,
		Image:       &image,
		Price:       &price,
		Rating:      uint8(0),
	}
}

// add new product by admin
func (pm *ProductModel) AddProduct(name, description, image string, price uint64) (int64, error) {
	//set ratings to 0 initiallu
	product := NewProduct(name, description, image, price)
	query := `insert into products(product_id, product_name, description, image, price, rating) values(?, ?, ?, ?, ?)`

	tx, err := pm.DB.Begin()
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(query)
	if err != nil {
		return 0, err
	}

	result, err := stmt.Exec(product.ProductID, product.ProductName, product.Description, product.Image, product.Price, product.Rating)
	if err != nil {
		return 0, err
	}
	ProductId, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	if err := tx.Commit(); err != nil {
		return 0, err
	}
	return ProductId, nil
}

// get product for other Operations by Id
func (pm *ProductModel) GetProduct(productID int) (*models.ResponseProduct, error) {
	query := `select product_name, description, price, rating,image from products where product_id = ?`

	tx, err := pm.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	row, err := stmt.Query(productID)
	if err != nil {
		return nil, err
	}
	defer row.Close()
	var Product models.ResponseProduct
	if row.Next() {
		err = row.Scan(&Product.ProductName, &Product.Description, Product.Price, Product.Image)
		if err != nil {
			return nil, err
		}
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return &Product, nil
}

// search for Product by name
func (pm *ProductModel) GetProductByName(name string) ([]*models.ResponseProduct, error) {
	query := `select product_name, description, price, rating,image from products where product_name like ?`

	tx, err := pm.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	//create statament
	stmt, err := tx.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(name)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var Products []*models.ResponseProduct
	for rows.Next() {
		var product *models.ResponseProduct
		err := rows.Scan(&product.ProductName, &product.Description, &product.Price, &product.Rating, &product.Image)
		if err != nil {
			return nil, err
		}
		Products = append(Products, product)
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return Products, nil
}

// out of units
func (pm *ProductModel) RemoveProduct(productID int) error {
	query := `delete from products where product_id = ?`

	tx, err := pm.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(productID)
	if err != nil {
		return err
	}
	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

// add product to user cart
func (pm *ProductModel) AddProductoCart(userID int) error {

	return nil
}

// cart operations
func (pm *ProductModel) GetUserCart(userID int) ([]*models.ResponseCartProducts, error) {
	query := `select product_name, price, rating, image, quantity, color, size from carts where user_id = ?`

	tx, err := pm.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var userCart []*models.ResponseCartProducts
	for rows.Next() {
		var userProducts *models.ResponseCartProducts
		err := rows.Scan(&userProducts.ProductName, &userProducts.Price, userProducts.Rating, userProducts.Quantity, userProducts.Color, userProducts.Size)
		if err != nil {
			return nil, err
		}
		userCart = append(userCart, userProducts)
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return userCart, nil
}

// AddProductToCart adds a product to the user's cart.
func (pm *ProductModel) AddProductToCart(userID, productID int, productName string, price, rating int, image string, quantity int, color, size string) error {
	// Assuming pm.DB is a valid *sql.DB connection

	// Create the SQL query with placeholders
	query := `
		INSERT INTO carts (user_id, product_id, product_name, price, rating, image, quantity, color, size)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	// Execute the SQL query with the provided parameters
	_, err := pm.DB.Exec(query, userID, productID, productName, price, rating, image, quantity, color, size)
	if err != nil {
		return err
	}

	return nil
}
