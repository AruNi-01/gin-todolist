package main

import (
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"go-web-project-todolist/dao"
	"go-web-project-todolist/models"
	"go-web-project-todolist/routers"
)

func main() {
	// 创建数据库
	// SQL: CREATE DATABASE goweb_todolist

	// 连接数据库
	err := dao.InitMySQL()
	// 连接数据库失败属于不可逆错误，直接panic
	if err != nil {
		panic(err)
	}
	defer dao.Close()

	// Model绑定
	dao.DB.AutoMigrate(&models.Todo{})

	// 注册路由
	route := routers.SetupRouter()
	err = route.Run("127.0.0.1:8080")
	if err != nil {
		return
	}
}
