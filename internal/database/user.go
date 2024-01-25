package database

import (
	"database/sql"
	"errors"
	"strings"

	"github.com/h3th-IV/mysticMerch/internal/models"
	"github.com/h3th-IV/mysticMerch/internal/utils"
	"golang.org/x/crypto/bcrypt"
)

type UserModel struct {
	DB *sql.DB
}

func NewUser(firstName, lastName, email, phoneNumber string) *models.User {
	uuid, _ := utils.GenerateUUID("user")
	return &models.User{
		UserID:      uuid,
		FirstName:   &firstName,
		LastName:    &lastName,
		Email:       &email,
		PhoneNumber: &phoneNumber,
	}
}

// create new user in dB
func (um *UserModel) InsertUser(fname, lname, password, email, phoneNumber, UserID string) error {
	user := NewUser(fname, lname, email, phoneNumber)
	passwordHash, err := utils.HashPassword(password)
	if err != nil {
		return err
	}
	query := `insert into users(user_id, first_name, last_name, email, phone_number, password_hash) values(?, ?, ?, ?, ?, ?)`

	tx, err := um.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	//statement
	stmt, err := tx.Prepare(query)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(user.UserID, user.FirstName, user.LastName, user.Email, user.PhoneNumber, passwordHash)
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

// GetUserby email(i.e when logged in)
func (pm *ProductModel) GetUserID(email string) (int, error) {
	query := `select id from users where email = ?`
	tx, err := pm.DB.Begin()
	if err != nil {
		return 0, nil
	}
	defer tx.Rollback()
	stmt, err := tx.Prepare(query)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()
	user := models.User{}
	rowErr := stmt.QueryRow(email).Scan(&user.ID)
	if rowErr != nil {
		if errors.Is(rowErr, sql.ErrNoRows) {
			return 0, rowErr
		}
	}
	return *user.ID, nil
}

// auth the user for login
func (um *UserModel) AuthenticateUser(email, password string) (*models.User, error) {
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
	err = stmt.QueryRow(email).Scan(&user.ID, &user.UserID, &user.FirstName, &user.LastName, &user.Email, &user.PhoneNumber, &user.PasswordHash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, utils.ErrInvalidCredentials
		} else {
			return nil, err
		}
	}

	err = bcrypt.CompareHashAndPassword(user.PasswordHash, []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
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
func (um *UserModel) RemoveUser(user_id int) error {
	query := `delete from users where user_id = ?`

	//use db pool
	tx, err := um.DB.Begin()
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

//func (um *UserModel) AddUserAddress(UserId int)
