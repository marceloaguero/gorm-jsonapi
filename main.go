package main

import (
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Base contains common columns for all tables.
type Base struct {
	ID        string `gorm:"size:36"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

// BeforeCreate will set a UUID rather than numeric ID.
func (base *Base) BeforeCreate(tx *gorm.DB) error {
	base.ID = uuid.New().String()
	return nil
}

// User is the model for the user table.
type User struct {
	Base
	SomeFlag bool `gorm:"not null;default:true"`
	Profile  Profile
}

// Profile is the model for the profile table.
type Profile struct {
	Base
	Name   string `gorm:"size:60;not null;"`
	UserID string `gorm:"size:36"`
}

func main() {
	dsn := "root:supersecret@tcp(127.0.0.1:3306)/test?charset=utf8&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		panic(err)
	}

	// db.AutoMigrate(&User{}, &Profile{})

	user := &User{SomeFlag: false}
	if db.Create(&user).Error != nil {
		log.Panic("Unable to create user.")
	}

	profile := &Profile{Name: "Marcelo Aguero", UserID: user.Base.ID}
	if db.Create(&profile).Error != nil {
		log.Panic("Unable to create profile.")
	}

	fetchedUser := &User{}
	// if db.Where("id = ?", profile.UserID).Preload("Profile").First(&fetchedUser).RecordNotFound() {
	db.Where("id = ?", profile.UserID).Preload("Profile").First(&fetchedUser)
	//		log.Panic("Unable to find created user.")
	//	}

	fmt.Printf("User: %+v\n", fetchedUser)
}
