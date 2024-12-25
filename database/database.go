// database/database.go
package database

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var DB *gorm.DB

func InitDB() {
	dsn := "root:040131@tcp(127.0.0.1:3306)/book_management?charset=utf8mb4&parseTime=True&loc=Local"
	var err error
	DB, err = gorm.Open("mysql", dsn)
	if err != nil {
		panic(fmt.Sprintf("连接数据库失败: %v", err))
	}
}
