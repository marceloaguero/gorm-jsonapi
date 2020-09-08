package main

import (
	"fmt"
	"log"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

// Base contains common columns for all tables.
type Base struct {
	ID        uuid.UUID  `gorm:"type:uuid;primary_key;"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"update_at"`
	DeletedAt *time.Time `sql:"index" json:"deleted_at"`
}

// BeforeCreate will set a UUID rather than numeric ID.
func (base *Base) BeforeCreate(scope *gorm.Scope) error {
	uuid, err := uuid.NewV4()
	if err != nil {
		return err
	}
	return scope.SetColumn("ID", uuid)
}

// User is the model for the user table.
type User struct {
	Base
	SomeFlag bool    `gorm:"column:some_flag;not null;default:true" json:"some_flag"`
	Profile  Profile `json:"profile"`
}

// Profile is the model for the profile table.
type Profile struct {
	Base
	Name   string    `gorm:"column:name;size:128;not null;" json:"name"`
	UserID uuid.UUID `gorm:"type:uuid;column:user_foreign_key;not null;" json:"-"`
}

func main() {
	db, err := gorm.Open("sqlite3", "test.db")
	if err != nil {
		panic(err)
	}

	db.LogMode(true)
	db.AutoMigrate(&User{}, &Profile{})

	user := &User{SomeFlag: false}
	if db.Create(&user).Error != nil {
		log.Panic("Unable to create user.")
	}

	profile := &Profile{Name: "New User", UserID: user.Base.ID}
	if db.Create(&profile).Error != nil {
		log.Panic("Unable to create profile.")
	}

	fetchedUser := &User{}
	if db.Where("id = ?", profile.UserID).Preload("Profile").First(&fetchedUser).RecordNotFound() {
		log.Panic("Unable to find created user.")
	}

	fmt.Printf("User: %+v\n", fetchedUser)
}
