package routers

import (
	"github.com/gin-gonic/gin"
	"go-web-project-todolist/controller"
)

func SetupRouter() *gin.Engine {
	route := gin.Default()
	// 告诉gin 模板文件引用的静态文件在哪儿
	route.Static("/static", "static")
	// 告诉gin 模板文件在哪儿
	route.LoadHTMLGlob("templates/*")

	// 主页
	route.GET("/", controller.IndexHandler)

	// v1 路由组
	v1Group := route.Group("v1")
	{
		// 添加
		v1Group.POST("/todo", controller.AddTodo)

		// 查看
		// 查看所有待办事项
		v1Group.GET("/todo", controller.QueryAllTodos)

		// 根据id查看某个待办事项
		v1Group.GET("/todo/:id", controller.QueryTodoById)

		// 根据id修改待办事项的状态
		v1Group.PUT("todo/:id", controller.UpdateTodoStatusById)

		// 根据id删除代办事项
		v1Group.DELETE("todo/:id", controller.DeleteTodoById)
	}
	return route
}
