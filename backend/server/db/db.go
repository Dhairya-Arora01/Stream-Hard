package db

import (
	"errors"
	"fmt"
	"os"

	"github.com/Dhairya-Arora01/StreamHard/server/auth"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

// InitDb initializes the Database.
func InitDB() {
	// Connecting to the db.
	var err error
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s", os.Getenv("HOST"), os.Getenv("USER"), os.Getenv("PASSWORD"), os.Getenv("DBNAME"), os.Getenv("PORT"))
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	// Migrating the Schema
	DB.AutoMigrate(&User{})

}

// CreateUser creates the User
func CreateUser(Name, Email, Password string) error {

	if Name == "" || len(Password) < 6 {
		return errors.New("unable to create user")
	}

	x := DB.Create(&User{
		Name:     Name,
		Email:    Email,
		Password: Password,
	})

	if x.Error != nil {
		return x.Error
	}

	return nil
}

// LoginUser returns the username, token and error(if any).
func LoginUser(Email, Password string) (string, string, error) {

	var user User
	DB.Where(map[string]interface{}{
		"email":    Email,
		"password": Password,
	}).First(&user)

	if user.ID == 0 {
		return "", "", errors.New("no user with given details")
	}

	tokenString := auth.GenerateToken(int(user.ID))

	return user.Name, tokenString, nil
}
