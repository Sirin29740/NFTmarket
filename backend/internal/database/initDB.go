package database

import (
	"NFTmarket/internal/user"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() *gorm.DB {
	dsn := "root:11220518@tcp(localhost:3306)/nftmarket?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println("failed to connect databse")
	}
	err = db.AutoMigrate(&user.User{})
	if err != nil {
		fmt.Println("failed to auto migrate User")
	}
	DB = db
	return db
}
func GetDB() *gorm.DB {
	return DB
}
