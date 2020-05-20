package config

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

func DbConn() (db *sql.DB) {
	//dbdriver := dbDriver
	//dbuser := dbUser
	//dbpass := dbPass
	//dbname := dbName
	dbdriver := "mysql"
	dbuser := "root"
	dbpass := "!Q2w#E4r"
	dbname := "api_go"
	db, err := sql.Open(dbdriver, dbuser+":"+dbpass+"@tcp(127.0.0.1:3306)/"+dbname)
	if err != nil {
		panic(err.Error())
	}
	return db
}
