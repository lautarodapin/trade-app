package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email" gorm:"unique;not null"`
	Token     string    `json:"token"`
	Password  string    `json:"password"`
	FavPairs  []FavPair `json:"fav_pairs"`
}

type Pair struct {
	gorm.Model
	Symbol   string    `json:"symbol" gorm:"unique;not null"`
	FavPairs []FavPair `json:"fav_pairs"`
}

type FavPair struct {
	gorm.Model
	UserID uint `json:"user_id" gorm:"not null"`
	User   User `json:"user" gorm:"foreignkey:UserID"`
	PairID uint `json:"pair_id" gorm:"not null"`
	Pair   Pair `json:"pair" gorm:"foreignkey:PairID"`
}

const (
	BUY  = 1
	SELL = 2
)

type Trade struct {
	gorm.Model
	UserID   uint    `json:"userId" gorm:"not null"`
	Type     uint8   `json:"type" gorm:"not null;default:1"`
	Quantity float64 `json:"quantity" gorm:"not null"`
	Price    float64 `json:"price" gorm:"not null"`
	Earns    float64 `json:"earns"`
}

func AutoMigrate(db *gorm.DB) {
	db.AutoMigrate(&User{})
	db.AutoMigrate(&Pair{})
	db.AutoMigrate(&FavPair{})
	db.AutoMigrate(&Trade{})
}

func InitUsers(db *gorm.DB) {
	// The email unique constraint will prevent from recreating the same users
	for _, user := range initialUsers {
		db.Create(&user)
	}
}

var initialUsers []User = []User{
	{FirstName: "John", LastName: "Doe", Email: "user1@test.com", Token: "d73e:9666:2dec:2ed8:073f:7f52:1ffc:5b9d", Password: HashPassword("password")},
	{FirstName: "Koko", LastName: "Doe", Email: "user2@test.com", Token: "9ae4:9c47:a59f:9427:bc36:f6ec:536f:3c83", Password: HashPassword("password")},
	{FirstName: "Francis", LastName: "Sunday", Email: "user3@test.com", Token: "f049:fc4e:eb2a:2d50:2962:5ab7:f5c7:6b96", Password: HashPassword("password")},
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
