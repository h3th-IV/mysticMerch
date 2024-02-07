// []TODO complete edit user adress
package database

import (
	"database/sql"
	"errors"
	"strings"

	"github.com/h3th-IV/mysticMerch/internal/models"
	"github.com/h3th-IV/mysticMerch/internal/utils"
)

// init new user
func NewUser(fName, LName, email, phoneNumber, password string) (*models.User, error) {
	crypted, err := utils.EncryptPass([]byte(password))
	if err != nil {
		return nil, err
	}

	uuid, err := utils.GenerateUUID("user")
	if err != nil {
		return nil, err
	}

	return &models.User{
		UserID:      uuid,
		FirstName:   &fName,
		LastName:    &LName,
		Email:       &email,
		PhoneNumber: &phoneNumber,
		Password:    crypted,
	}, nil
}

// create new user in dB
func (dm *DBModel) InsertUser(fname, lname, email, phoneNumber, password string) error {
	user, err := NewUser(fname, lname, email, phoneNumber, password)
	if err != nil {
		return err
	}
	query := `insert into users(user_id, first_name, last_name, email, phone_number, password_hash) values(?, ?, ?, ?, ?, ?)`

	tx, err := dm.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	//statement
	stmt, err := tx.Prepare(query)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(user.UserID, user.FirstName, user.LastName, user.Email, user.PhoneNumber, user.Password)
	if err != nil {
		//check if err is of type mysql err
		if errors.As(err, &utils.MySQLErr) {
			//check if error is existing credentials (not unique) with the constraint 'users_uc_email'
			if utils.MySQLErr.Number == 1062 && strings.Contains(utils.MySQLErr.Message, "user_uc_email") {
				return utils.ErrInvalidCredentials
			}
		}
		return err
	}
	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

// GetUserby uuid(i.e when logged in)
func (dm *DBModel) GetUserbyUUID(uuid string) (*models.ResponseUser, error) {
	query := `select id, first_name, last_name, email, phone_number from users where user_id = ?`
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
	user := models.ResponseUser{}
	rowErr := stmt.QueryRow(uuid).Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.PhoneNumber)
	if rowErr != nil {
		if errors.Is(rowErr, sql.ErrNoRows) {
			return nil, rowErr
		}
	}
	return &user, nil
}

// Get al users from DB
func (dm *DBModel) GetAllUsers() ([]*models.ResponseUser, error) {
	query := `select id, first_name, last_name, email, phone_number from users`
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
	defer rows.Close()
	var Users []*models.ResponseUser
	for rows.Next() {
		var uSer *models.ResponseUser
		rows.Scan(&uSer.ID, &uSer.FirstName, &uSer.LastName, &uSer.PhoneNumber)
		if err != nil {
			return nil, err
		}
		Users = append(Users, uSer)
	}
	return Users, nil
}

// auth the user for login
func (um *DBModel) AuthenticateUser(email, password string) (*models.User, error) {
	query := `select id, password_hash from users where email = ?`

	//use transaction
	tx, err := um.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	user := models.User{}
	err = stmt.QueryRow(email).Scan(&user.ID, &user.UserID, &user.FirstName, &user.LastName, &user.Email, &user.PhoneNumber, &user.Password)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, utils.ErrInvalidCredentials
		} else {
			return nil, err
		}
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Remove user fields --cascade set to on delete
func (dm *DBModel) RemoveUser(user_id int) error {
	query := `delete from users where user_id = ?`

	//use db pool
	tx, err := dm.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(query)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(user_id)
	if err != nil {
		return err
	}
	return nil
}

// create new address
func NewAddress(user *models.ResponseUser, houseNo, str, city, postalCode string) *models.Address {
	return &models.Address{
		HouseNo:     &houseNo,
		Street:      &str,
		City:        &city,
		PostalCode:  &postalCode,
		UserPhoneNo: &user.PhoneNumber,
	}
}

// register new address
func (dm *DBModel) AddUserAddress(user *models.ResponseUser, houseNo, str, city, postalCode string) error {
	adrr := NewAddress(user, houseNo, str, city, postalCode)
	query := `insert into address(user_id, house_no, street, city, postal_code) values(?, ?, ?, ?, ?)`

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
	_, err = stmt.Exec(user.ID, adrr.HouseNo, adrr.Street, adrr.City, adrr.PostalCode)
	if err != nil {
		return err
	}
	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

// return user addresses
func (dm *DBModel) ReturnUserAddress(userID int) ([]*models.Address, error) {
	query := `select * from address where user_id = ?`
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
	var UserAddrs []*models.Address
	for rows.Next() {
		var useraddr *models.Address
		if err := rows.Scan(&useraddr.HouseNo, &useraddr.Street, &useraddr.City, &useraddr.PostalCode); err != nil {
			return nil, err
		}
		UserAddrs = append(UserAddrs, useraddr)
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return UserAddrs, nil
}

// reomve address
func (dm *DBModel) ReomveAddress(userID, address_id int) error {
	query := `delete * from address where address_id = ? and user_id = ? `

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
	_, err = stmt.Exec(address_id, userID)
	if err != nil {
		return err
	}
	if err = tx.Commit(); err != nil {
		return err
	}
	return nil
}
