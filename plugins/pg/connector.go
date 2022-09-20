package pg

import (
	"animeGAN2go/plugins"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var db *gorm.DB

func init() {

	dbUri := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", plugins.DbHost,
		plugins.DbPort, plugins.Username, plugins.DbName, plugins.Password)

	fmt.Println(dbUri)

	conn, err := gorm.Open("postgres", dbUri)
	if err != nil {
		fmt.Print(err)
	}

	db = conn

}

//func GetDB() *gorm.DB {
//	return db
//}
