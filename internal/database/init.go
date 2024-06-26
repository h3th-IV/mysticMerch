package database

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/h3th-IV/mysticMerch/internal/utils"
)

type DBModel struct {
	DB *sql.DB
}

func InitDB() (*sql.DB, error) {
	if err := utils.LoadEnv(); err != nil {
		return nil, err
	}
	//logger package
	logger := utils.NewLogger(os.Stdout, os.Stderr)
	constr := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", os.Getenv("MM_USER"), os.Getenv("MM_PASSWORD"), os.Getenv("MM_HOST"), os.Getenv("MM_PORT"), os.Getenv("MM_DBNAME"))

	//open databse pool
	database, err := sql.Open("mysql", constr)
	if err != nil {
		logger.ErrLogger.Fatal(err)
		os.Exit(1)
	}

	//try database connection
	err = database.Ping()
	if err != nil {
		database.Close()
		logger.InfoLogger.Println("err connecting to database")
		os.Exit(1)
	}

	logger.InfoLogger.Println("connected to database succesfully")

	return database, nil
}

// CloseDB function
func (um *DBModel) CloseDB() error {
	if um.DB != nil {
		err := um.DB.Close()
		if err != nil {
			return err
		}
	}
	fmt.Println("Connection to Database Closed succesfully")
	return nil
}
