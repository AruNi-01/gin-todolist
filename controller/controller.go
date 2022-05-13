package controller

import (
	"github.com/gin-gonic/gin"
	"go-web-project-todolist/models"
	"net/http"
)

// IndexHandler 主页
func IndexHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", nil)
}

// AddTodo 添加代办事项
func AddTodo(c *gin.Context) {
	// 前面页面填写代办事项，点击提交，会发请求到这里
	// 1. 从请求中把数据拿出来
	var todo models.Todo
	c.BindJSON(&todo)

	// 2.插入数据库，返回响应
	// 插入数据库失败，状态码为200表示请求成功，把插入失败的error返回
	if err := models.AddTodo(&todo); err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, todo) // 插入成功将该数据返回
		
		/*
			企业中返回信息可能如下：
			c.JSON(http.StatusOK, gin.H{
				"code": 200,
				"msg": "success",
				"data": todo,
			})
		*/
	}
}

// QueryAllTodos 查看所有待办事项
func QueryAllTodos(c *gin.Context) {
	// 从数据库中查询出所有todolist，返回响应
	if todoList, err := models.QueryAllTodos(); err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, todoList)
	}
}

// QueryTodoById 根据id查看某个待办事项
func QueryTodoById(c *gin.Context) {
	// 1. 获取前端传来的参数
	id, ok := c.Params.Get("id")
	if !ok {
		c.JSON(http.StatusOK, gin.H{"error": "接收id失败"})
		return
	}

	// 2. 去数据库中查数据，返回响应
	if todo, err := models.QueryTodoById(id); err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, todo)
	}
}

// UpdateTodoStatusById 根据id修改待办事项的状态
func UpdateTodoStatusById(c *gin.Context) {
	// 获取前端传来的id
	id, ok := c.Params.Get("id")
	if !ok {
		c.JSON(http.StatusOK, gin.H{"error": "接收id失败"})
		return
	}

	// 再根据当前todo的status去修改status
	if err := models.UpdateTodoStatusById(id); err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
	} else {
		// 修改成功返回响应
		c.JSON(http.StatusOK, gin.H{id: "修改成功"})
	}
}

// DeleteTodoById 根据id删除代办事项
func DeleteTodoById(c *gin.Context) {
	// 获取前端传来的id
	id, ok := c.Params.Get("id")
	if !ok {
		c.JSON(http.StatusOK, gin.H{"error": "接收id失败"})
		return
	}

	// 删除
	if err := models.DeleteTodoById(id); err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{id: "delete成功"})
	}
}
