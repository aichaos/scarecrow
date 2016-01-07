/*
Package models defines database models for Scarecrow.
*/
package models

import (
	"time"
	"github.com/jinzhu/gorm"
)

type AdminSetting struct {
	// Whether users can register their own accounts
	OpenRegistration bool `sql:"DEFAULT:false"`
}

type User struct {
	ID          int
	Username    string
	Password    string
	Role        string `sql:"size:10;default:'user'"`
	AccessToken string // Personal API access token
	Created     time.Time
	Updated     time.Time

	Bots []Bot
}

type Bot struct {
	ID        string
	UserID    int
	Name      string
	BrainType string `sql:"DEFAULT:'RiveScript'"`
	Base      string // Refers to a Bot.ID
	IsBase    bool   `sql:"DEFAULT:false"`

	Listeners []Listener
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
