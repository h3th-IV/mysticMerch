package database

import (
	"database/sql"

	"github.com/h3th-IV/mysticMerch/internal/models"
)

type ProductModel struct {
	DB *sql.DB
}

// add new product by admin
func (pm *ProductModel) AddProduct(name, description, price, image string) (int64, error) {
	//set ratings to 0 initiallu
	rating := 0
	query := `insert into products(product_name, description, price, rating, image) values(?, ?, ?, ?, ?)`

	tx, err := pm.DB.Begin()
	if err != nil {
		return 0, err
	}

	stmt, err := tx.Prepare(query)
	if err != nil {
		return 0, err
	}

	result, err := stmt.Exec(name, description, price, rating, image)
	if err != nil {
		return 0, err
	}
	ProductId, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return ProductId, nil
}

// get product for other Operations by Id
func (pm *ProductModel) GetProduct(productID int) (*models.Product, error) {
	query := `select * from products where product_id = ?`

	result, err := pm.DB.Query(query, productID)
	if err != nil {
		return nil, err
	}
	Product := models.Product{}
	err = result.Scan(&Product.ProductID, &Product.ProductName, &Product.Description, Product.Price, Product.Image)
	if err != nil {
		return nil, err
	}
	return &Product, nil
}

// search for Product by name
func (pm *ProductModel) GetProductByName(name string) ([]*models.Product, error) {
	query := `select * from products where product_name like ?`
	rows, err := pm.DB.Query(query, name)
	if err != nil {
		return nil, err
	}

	var Products []*models.Product
	for rows.Next() {
		var product *models.Product
		err := rows.Scan(&product.ProductID, &product.ProductName, &product.Description, &product.Price, &product.Rating, &product.Image)
		if err != nil {
			return nil, err
		}
		Products = append(Products, product)
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

	_, err = stmt.Exec(productID)
	if err != nil {
		return err
	}
	return nil
}

// cart operations
func (pm *ProductModel) GetUserCart(userID int) ([]*models.UserProducts, error) {
	query := `select * from carts where user_id = ?`

	tx, err := pm.DB.Begin()
	if err != nil {
		return nil, err
	}
	stmt, err := tx.Prepare(query)
	if err != nil {
		return nil, err
	}
	rows, err := stmt.Query(userID)
	if err != nil {
		return nil, err
	}

	for rows.Next() {

	}

}
