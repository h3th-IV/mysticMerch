package models

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/h3th-IV/mysticMerch/internal/utils"
)

func InitDB() (*sql.DB, error) {
	//logger package
	logger := utils.NewLogger(os.Stdout, os.Stderr)
	constr := ""

	//open databse pool
	database, err := sql.Open("mysql", constr)
	if err != nil {
		logger.ErrLogger.Fatal(err)
		return nil, fmt.Errorf("err creating database pool: %v", err)
	}

	//try database connection
	err = database.Ping()
	if err != nil {
		database.Close()
		return nil, fmt.Errorf("err connecting to database: %v", err)
	}

	logger.InfoLogger.Println("connected to Database was succesfull")

	return database, nil
}

//closeDB function
