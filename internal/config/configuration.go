package config

import (
	"database/sql"
	"fmt"
)

func DBConnect() *sql.DB {
	dbDriver := "mysql"
	dbUser := "developer"
	dbPassword := "developer"
	dbUrl := "127.0.0.1:3306"
	dbName := "CRUD"
	URL := fmt.Sprintf("%s:%s@tcp(%s)/%s", dbUser, dbPassword, dbUrl, dbName)
	db, err := sql.Open(dbDriver, URL)
	if err != nil {
		panic(err.Error())
	} else {
		err = db.Ping()
		if err != nil {
			panic(err.Error())
		} else {
			fmt.Println("Connected to Database !")
		}
	}

	return db
}
