package models

import (
	"log"
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectToDB() {
	var err error
	user := "root"
	password := "password"
	hostname := "127.0.0.1"
	port := "3306"
	dbname := "pointsdb"
	
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", user, password, hostname, port, dbname) + "?charset=utf8mb4&parseTime=True&loc=Local"
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("Failed to connect to Database")
	}
}