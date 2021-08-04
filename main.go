package main

import (
	"crud-golang/internal/config"
	"crud-golang/internal/routes"
	"crud-golang/pkg/database"

	_ "github.com/go-sql-driver/mysql"
)

func main() {

	db := config.DBConnect()
	database.SetDB(db)
	defer db.Close()

	routes.INIT_Routes()

}
