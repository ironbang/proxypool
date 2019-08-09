package database

import (
	"fmt"
	"github.com/ironbang/proxypool/common/config"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"log"
)

var db *gorm.DB

func init() {
	var err error
	user := config.GetDatabase("User")
	password := config.GetDatabase("Password")
	ip := config.GetDatabase("IP")
	port := config.GetDatabase("Port")
	tablename := config.GetDatabase("Tablename")
	args := fmt.Sprintf("%s:%s@(%s:%s)/%s?charset=utf8&parseTime=True", user, password, ip, port, tablename)
	if db, err = gorm.Open("mysql", args); err != nil {
		log.Fatal(err.Error())
	}
	db.DB().SetMaxOpenConns(10)
	fmt.Println("连接数据库成功")
}

func GetDB() *gorm.DB {
	return db
}
