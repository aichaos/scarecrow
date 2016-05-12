/*
Package models defines database models for Scarecrow.
*/
package models

import (
	"github.com/jinzhu/gorm"
)

type AppConfig struct {
	Initialized bool   `gorm:"primary_key"`
	Username    string // Admin username
	Password    string // Admin password
	Name        string // Bot name
	Replies     string // RiveScript reply root
}

type Listener struct {
	ID       int
	Type     string
	Enabled  bool
	Settings []ListenerSetting
}

type ListenerSetting struct {
	gorm.Model
	ListenerID int
	Key        string
	Value      string
}

type User struct {
	ID   string
	Data []UserData
}

type UserData struct {
	gorm.Model
	UserID string
	Key    string
	Value  string
}
