package database

import (
	"github.com/h3th-IV/mysticMerch/internal/models"
	"github.com/h3th-IV/mysticMerch/internal/utils"
)

// TODO: push this to the products package.

func NewProduct(name, description, image string, price uint64) (*models.Product, error) {
	uuid, err := utils.GenerateUUID("product")
	return &models.Product{
		ProductID:   &uuid,
		ProductName: &name,
		Description: &description,
		Image:       &image,
		Price:       &price,
		Rating:      uint8(0),
	}, err
}

/* admin operations*/

// add new product by admin
func (dm *DBModel) AddProduct(name, description, image string, price uint64) (int64, error) {
	//set ratings to 0 initiallu
	product, err := NewProduct(name, description, image, price)
	if err != nil {
		return 0, err
	}
	query := `insert into products(product_id, product_name, description, image, price, rating) values(?, ?, ?, ?, ?)`

	tx, err := dm.DB.Begin()
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

// out of units  --admin stuff
func (dm *DBModel) RemoveProduct(productID int) error {
	query := `delete from products where product_id = ?`

	tx, err := dm.DB.Begin()
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

/* Normal Product operations */

// viewProducts --a list of products for home page
func (dm *DBModel) ViewProducts() ([]*models.ResponseProduct, error) {
	query := `select top 30 product_name, description, image, price, rating from products`

	tx, err := dm.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		return nil, err
	}
	var Products []*models.ResponseProduct
	for rows.Next() {
		var product *models.ResponseProduct
		rows.Scan(&product.ProductName, &product.Description, &product.Image, &product.Price, &product.Rating)
		if err != nil {
			return nil, err
		}
		Products = append(Products, product)
	}
	return Products, nil
}

// get product for other Operations by product uuid
func (dm *DBModel) GetProduct(productID string) (*models.Product, error) {
	query := `select * from products where product_id = ?`

	tx, err := dm.DB.Begin()
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
	var Product models.Product
	if row.Next() {
		err = row.Scan(&Product.ProductID, &Product.ProductName, &Product.Description, &Product.Price, &Product.Image)
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
func (dm *DBModel) GetProductByName(name string) ([]*models.ResponseProduct, error) {
	query := `select product_name, description, price, rating,image from products where product_name like ?`

	tx, err := dm.DB.Begin()
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
