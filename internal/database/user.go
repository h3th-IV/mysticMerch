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

// create new user in dB
func (um *UserModel) InsertUser(fname, lname, password, email, phoneNumber, UserID string) error {
	passwordHash, err := utils.HashPassword(password)
	if err != nil {
		return err
	}
	query := `insert into users(first_name, last_name, password_hash, email, phone_number, user_id) values(?, ?, ?, ?, ?, ?)`

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

	_, err = stmt.Exec(fname, lname, passwordHash, email, phoneNumber, UserID)
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

// auth the user for login
func (um *UserModel) AuthenticateUser(email, password string) (int, error) {
	query := `select id, password_hash from users where email = ?`

	//use transaction
	tx, err := um.DB.Begin()
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(query)
	if err != nil {
		return 0, err
	}
	user := models.User{}
	err = stmt.QueryRow(email).Scan(&user.ID, &user.PasswordHash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, utils.ErrInvalidCredentials
		} else {
			return 0, err
		}
	}

	err = bcrypt.CompareHashAndPassword(user.PasswordHash, []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return 0, utils.ErrInvalidCredentials
		} else {
			return 0, err
		}
	}
	err = tx.Commit()
	if err != nil {
		return 0, err
	}
	return *user.ID, nil
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
