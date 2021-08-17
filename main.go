package main

import (
	"crud-golang/internal/config"
	"crud-golang/internal/routes"
	"crud-golang/pkg/database"

	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

func main() {

	db := config.DBConnect()
	database.SetDB(db)
	defer db.Close()

	tracer.Start()
	defer tracer.Stop()

	routes.INIT_Routes()
}
