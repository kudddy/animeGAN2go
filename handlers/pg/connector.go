package handlers

import (
	handlers "animeGAN2go/handlers/plugins"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var db *gorm.DB

// pq: unsupported sslmode "enable"; only "require" (default), "verify-full", "verify-ca", and "disable" supportedget file ids from queen
func init() {

	dbUri := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=require password=%s", plugins.DbHost,
		handlers.DbPort, handlers.Username, handlers.DbName, handlers.Password)

	fmt.Println(dbUri)

	conn, err := gorm.Open("postgres", dbUri)
	if err != nil {
		fmt.Println("troubles with connect to database")
		fmt.Print(err)
	}

	db = conn

}

func GetDB() *gorm.DB {
	return db
}
