package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	userID   uint   `gorm:"primaryKey"`
	username string `gorm:"unique"`
	password string
	email    string `gorm:"unique"`
	polls    []Poll `gorm:"one2many"`
}

type Poll struct {
	gorm.Model
	pollID  uint     `gorm:"primaryKey"`
	title   string   `gorm:"not null"`
	userID  User     `gorm:"foreignKey:UserID"`
	options []Option `gorm:"one2many"`
}

type Option struct {
	gorm.Model
	optionID uint   `gorm:"primaryKey"`
	content  string `gorm:"not null"`
	pollID   Poll   `gorm:"foreignKey:pollID"`
	correct  bool
}
