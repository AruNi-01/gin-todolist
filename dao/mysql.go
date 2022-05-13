package dao

import "github.com/jinzhu/gorm"
import _ "github.com/jinzhu/gorm/dialects/mysql"

// DB 定义全局变量DB
var (
	DB *gorm.DB
)

// InitMySQL 连接数据库，成功返回nil，失败返回error
func InitMySQL() (err error) {
	dsn := "root:123456@tcp(127.0.0.1:3306)/goweb_todolist?charset=utf8mb4&parseTime=True&loc=Local"

	// 初始化全局变量DB
	DB, err = gorm.Open("mysql", dsn)
	if err != nil {
		return
	}
	return DB.DB().Ping() // 测试连通性再返回,ping得通返回nil，否则返回error
}

func Close() {
	DB.Close()
}
