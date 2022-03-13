package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	Email     string    `json:"email" gorm:"unique;not null"`
	Token     string    `json:"token"`
	Password  string    `json:"password"`
	FavPairs  []FavPair `json:"favPairs"`
}

type Pair struct {
	gorm.Model
	Symbol   string    `json:"symbol" gorm:"unique;not null"`
	FavPairs []FavPair `json:"favPairs"`
}

type FavPair struct {
	gorm.Model
	UserID uint `json:"userId" gorm:"not null"`
	PairID uint `json:"pairId" gorm:"not null"`
}

func AutoMigrate(db *gorm.DB) {
	db.AutoMigrate(&User{})
	db.AutoMigrate(&Pair{})
	db.AutoMigrate(&FavPair{})
}

func InitUsers(db *gorm.DB) {
	// The email unique constraint will prevent from recreating the same users
	for _, user := range initialUsers {
		db.Create(&user)
	}
}

var initialUsers []User = []User{
	{FirstName: "John", LastName: "Doe", Email: "user1@test.com", Token: "d73e:9666:2dec:2ed8:073f:7f52:1ffc:5b9d", Password: "password"},
	{FirstName: "Koko", LastName: "Doe", Email: "user2@test.com", Token: "9ae4:9c47:a59f:9427:bc36:f6ec:536f:3c83", Password: "password"},
	{FirstName: "Francis", LastName: "Sunday", Email: "user3@test.com", Token: "f049:fc4e:eb2a:2d50:2962:5ab7:f5c7:6b96", Password: "password"},
}
var initialPartList []string = []string{
	"BTCUSDT",
	"ETHUSDT",
	"BNBUSDT",
	"BCCUSDT",
	"NEOUSDT",
	"LTCUSDT",
	"QTUMUSDT",
	"ADAUSDT",
	"XRPUSDT",
	"EOSUSDT",
}

func InitPairList(db *gorm.DB) {
	db.Unscoped().Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&Pair{})
	for _, pair := range initialPartList {
		db.Create(&Pair{Symbol: pair})
	}

}

func InitFavPairList(db *gorm.DB, addInitialValues bool) {
	db.Unscoped().Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&FavPair{})
	var users []User
	if addInitialValues {
		var favPairs []FavPair
		db.Find(&users)
		for _, user := range users {
			db.Find(&favPairs).Where("user_id = ?", user.ID)
			for _, initialPair := range initialPartList[:3] {
				var pair Pair
				db.Where("symbol = ?", initialPair).First(&pair)
				db.Create(&FavPair{UserID: user.ID, PairID: pair.ID})
			}
		}
	}
}
