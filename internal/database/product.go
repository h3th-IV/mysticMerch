package database

import (
	"errors"

	"github.com/h3th-IV/mysticMerch/internal/models"
	"github.com/h3th-IV/mysticMerch/internal/utils"
)

// TODO: push this to the products package.

func NewProduct(name, description, image string, price float64) (*models.Product, error) {
	uuid, err := utils.GenerateUUID("product")
	return &models.Product{
		ProductID:   uuid,
		ProductName: name,
		Description: description,
		Image:       image,
		Price:       price,
		Rating:      int8(0),
	}, err
}

/* admin operations*/

// add new product by admin
func (dm *DBModel) AddProduct(adminID int, name, description, image string, price float64) (int64, error) {
	//set ratings to 0 initially
	if adminID != 1 {
		return 0, errors.New("only admin can add products")
	}
	product, err := NewProduct(name, description, image, price)
	if err != nil {
		return 0, err
	}
	query := `insert into products(product_id, product_name, description, image, price, rating) values(?, ?, ?, ?, ?, ?)`

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

// check for product existence
func (dm *DBModel) CheckProductExist(productUUID string) (int, error) {
	query := `select exists (select 1 from products where product_id = ?) as product_exists`

	tx, err := dm.DB.Begin()
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(query)
	if err != nil {
		return 0, err
	}
	var productExists int
	err = stmt.QueryRow(productUUID).Scan(&productExists)
	if err != nil {
		return 0, err
	}
	return productExists, nil
}

// out of units  --admin stuff
func (dm *DBModel) RemoveProductFromStore(productUUID string) error {
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

	_, err = stmt.Exec(productUUID)
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
func (dm *DBModel) ViewHomeProducts() ([]*models.ResponseProduct, error) {
	query := `select product_name, description, image, price, rating from products limit 30`

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
		product := &models.ResponseProduct{}
		if err := rows.Scan(&product.ProductName, &product.Description, &product.Image, &product.Price, &product.Rating); err != nil {
			return nil, err
		}
		if err := rows.Err(); err != nil {
			return nil, err
		}
		Products = append(Products, product)
	}
	return Products, nil
}

// get product for other Operations by product uuid
func (dm *DBModel) GetProduct(productUUID string) (*models.Product, error) {
	query := `select product_name, description, image, price, rating from products where product_id = ?`

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

	row, err := stmt.Query(productUUID)
	if err != nil {
		return nil, err
	}
	defer row.Close()
	var Product models.Product
	if row.Next() {
		err = row.Scan(&Product.ProductName, &Product.Description, &Product.Image, &Product.Price, &Product.Rating)
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
		product := &models.ResponseProduct{}
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
