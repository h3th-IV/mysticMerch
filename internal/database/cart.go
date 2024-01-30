package database

import (
	"fmt"

	"github.com/h3th-IV/mysticMerch/internal/models"
)

/* cart operations */

// add product to user cart

// view user cart
func (dm *DBModel) GetUserCart(userID int) ([]*models.ResponseCartProducts, error) {
	query := `select product_name, price, rating, image, quantity, color, size from carts where user_id = ?`

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

	rows, err := stmt.Query(userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var userCart []*models.ResponseCartProducts
	for rows.Next() {
		var userProducts *models.ResponseCartProducts
		err := rows.Scan(&userProducts.ProductName, &userProducts.Price, &userProducts.Rating, &userProducts.Quantity, &userProducts.Color, userProducts.Size)
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

// add product to user cart
func (dm *DBModel) AddProductoCart(userID, quantity int, productID, color, size string) error {
	query := `insert into carts(user_id, product_id, product_name, price, rating, image, quantity, color, size) values(?, ?, ?, ?, ?, ?, ?, ?, ?)`
	//retrive product info
	product, err := dm.GetProduct(productID)
	if err != nil {
		return err
	}
	if product == nil {
		return fmt.Errorf("product with uuid, %v not found", productID)
	}
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
	_, err = stmt.Exec(userID, product.ProductID, product.ProductName, product.Price, product.Rating, product.Image, quantity, color, size)
	if err != nil {
		return err
	}
	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

// check for product in user cart
func (dm *DBModel) CheckProductExistInUserCart(userid int, productId string) (bool, error) {
	query := `select count(*) from products where user_id = ? and product_id = ?`
	tx, err := dm.DB.Begin()
	if err != nil {
		return false, err
	}
	defer tx.Rollback()
	stmt, err := tx.Prepare(query)
	if err != nil {
		return false, err
	}
	defer stmt.Close()
	var count int
	countErr := stmt.QueryRow(userid, productId).Scan(&count)
	if countErr != nil {
		return false, countErr
	}
	return count > 0, nil
}

// pop item from user cart
func (dm *DBModel) RemoveItemfromCart(userid int, productID string) error {
	query := `delete from carts where user_id = ? and product = ?`
	productChecker, err := dm.CheckProductExistInUserCart(userid, productID)
	if err != nil {
		return err
	}
	if !productChecker {
		return fmt.Errorf("product with productID %v does not exist in user's cart", productID)
	}

	tx, err := dm.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	stmt, err := tx.Prepare(query)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(userid, productID)
	if err != nil {
		return err
	}
	if err = tx.Commit(); err != nil {
		return err
	}
	return nil
}

// edit cart like size, color, quantity
