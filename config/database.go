package config

import (
	"fmt"

	"github.com/HaikalRFadhilahh/course-golang/helper"
	"github.com/HaikalRFadhilahh/course-golang/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func DatabaseConnection() (*gorm.DB, error) {
	// Setting Database Connection Variable
	var (
		host     = helper.GetEnv("DB_HOST", "127.0.0.1")
		port     = helper.GetEnv("DB_PORT", "3306")
		user     = helper.GetEnv("DB_USERNAME", "root")
		password = helper.GetEnv("DB_PASSWORD", "")
		dbName   = helper.GetEnv("DB_NAME", "")
	)

	// Create Connection String
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", user, password, host, port, dbName)

	// Create Database Connection
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Database Refused!")
	}

	// Return Value if Database Success Create
	return db, err
}

func CreateOwnerAccount(db *gorm.DB) {
	hashPassword, err := bcrypt.GenerateFromPassword([]byte("123456"), bcrypt.DefaultCost)
	if err != nil {
		panic(fmt.Sprintf("Error Migration : %s", err.Error()))
	}

	owner := models.User{
		Role:     "Admin",
		Name:     "Owner",
		Password: string(hashPassword),
		Email:    "owner@go.id",
	}

	if db.Where("email=?", owner.Email).First(&owner).RowsAffected == 0 {
		db.Create(&owner)
	} else {
		fmt.Println("Migration not running,Data has exist")
	}
}
