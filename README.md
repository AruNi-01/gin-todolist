# gin-todolist
Gin框架小练习

## 前端

在GitHub中拉取前端的代码：https://github.com/Q1mi/bubble_frontend

在本地IDE中打开，构建项目：

```bash
npm install
npm run build
```

构建完后将dist目录中的文件拷贝到后端项目中，供后端使用。



## 后端

初始目录结构：

![image-20220512173603618](https://run-notes.oss-cn-beijing.aliyuncs.com/notes/992c63b69a9177567c7b4ca3f4bfbdee.png)



下载安装Gin：

```bash
go get -u github.com/gin-gonic/gin
```

编写`main.go`：

```go
package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	route := gin.Default()
	// 告诉gin 模板文件引用的静态文件在哪儿
	route.Static("/static", "static")
	// 告诉gin 模板文件在哪儿
	route.LoadHTMLGlob("templates/*")

	// 主页
	route.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	err = route.Run()
	if err != nil {
		return 
	}
}
```

运行项目，访问`localhost:8080/`：

![image-20220512174242212](https://run-notes.oss-cn-beijing.aliyuncs.com/notes/7ebe2b324f7841a8c29d15e0e715da42.png)

项目运行成功，下面开始编写核心代码。



## 框架搭建

创建一个名为`goweb_todolist`的MySQL数据库，然后引入GORM框架进行初始化：

```go
package main

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"net/http"
)

// DB 定义全局变量DB
var (
	DB *gorm.DB
)

// Todo Model：代办清单结构体
type Todo struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Status bool   `json:"status"`
}

// 连接数据库，成功返回nil，失败返回error
func initMySQL() (err error) {
	dsn := "root:123456@tcp(127.0.0.1:3306)/goweb_todolist?charset=utf8mb4&parseTime=True&loc=Local"

	// 初始化全局变量DB
	DB, err = gorm.Open("mysql", dsn)
	if err != nil {
		return
	}
	return DB.DB().Ping() // 测试连通性再返回,ping得通返回nil，否则返回error
}

func main() {
	// 创建数据库
	// SQL: CREATE DATABASE goweb_todolist

	// GORM框架---------------------
	// 连接数据库
	err := initMySQL()
	// 连接数据库失败属于不可逆错误，直接panic
	if err != nil {
		panic(err)
	}
	defer DB.Close()

	// Model绑定
	DB.AutoMigrate(&Todo{})

	// Gin框架---------------------
	route := gin.Default()
	// 告诉gin 模板文件引用的静态文件在哪儿
	route.Static("/static", "static")
	// 告诉gin 模板文件在哪儿
	route.LoadHTMLGlob("templates/*")

	// 主页
	route.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	// v1 路由组
	v1Group := route.Group("v1")
	{
		// 代办事项
		// 添加
		v1Group.POST("/todo", func(c *gin.Context) {

		})

		// 查看
		// 查看所有待办事项
		v1Group.GET("/todo", func(c *gin.Context) {

		})
		// 查看某个待办事项
		v1Group.GET("/todo/:id", func(c *gin.Context) {

		})

		// 修改
		v1Group.PUT("todo/:id", func(c *gin.Context) {

		})

		// 删除
		v1Group.DELETE("todo/:id", func(c *gin.Context) {

		})
	}

	err = route.Run()
	if err != nil {
		return 
	}
}
```

框架搭建好了，下面的具体实现就简单了。



## 功能实现

### 添加待办

```go
		// 添加
		v1Group.POST("/todo", func(c *gin.Context) {
			// 前面页面填写代办事项，点击提交，会发请求到这里
			// 1. 从请求中把数据拿出来
			var todo Todo
			c.BindJSON(&todo)

			// 2. 插入数据库
			err = DB.Create(&todo).Error

			// 3. 返回响应
			// 插入数据库失败，状态码为200表示请求成功，把插入失败的error返回
			if err != nil {
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
		})
```



### 查看代办

#### 查看所有代办

```go
		// 查看所有待办事项
		v1Group.GET("/todo", func(c *gin.Context) {
			// 1. 从数据库中查询出所有todolist
			var todoList []Todo
			err = DB.Find(&todoList).Error
			// 2. 返回响应
			if err != nil {
				c.JSON(http.StatusOK, gin.H{"error": err.Error()})
			} else {
				c.JSON(http.StatusOK, todoList)
			}
		})
```



#### 查看某个代办

```go
		// 查看某个待办事项
		v1Group.GET("/todo/:id", func(c *gin.Context) {
			// 1. 获取前端传来的参数
			id, ok := c.Params.Get("id")
			if !ok {
				c.JSON(http.StatusOK, gin.H{"error": "接收id失败"})
				return
			}
			// 2. 去数据库中查数据
			var todo Todo
			err = DB.First(&todo, id).Error
			// 3. 返回响应
			if err != nil {
				c.JSON(http.StatusOK, gin.H{"error": err.Error()})
			} else {
				c.JSON(http.StatusOK, todo)
			}
		})
```



### 修改代办

```go
		// 修改
		v1Group.PUT("todo/:id", func(c *gin.Context) {
			// 获取前端传来的id
			id, ok := c.Params.Get("id")
			if !ok {
				c.JSON(http.StatusOK, gin.H{"error": "接收id失败"})
				return
			}
			// 先查询出当前的todo事项，才知道修改成什么
			var todo Todo
			err = DB.Where("id = ?", id).First(&todo).Error

			// 再根据当前todo的status去修改status
			if err != nil {
				c.JSON(http.StatusOK, gin.H{"error": err.Error()})
			} else {
				if todo.Status {
					DB.Model(&todo).Select("status").Update(map[string]interface{}{"status": false})
				} else {
					DB.Model(&todo).Select("status").Update(map[string]interface{}{"status": true})
				}
				// 修改完后返回响应
				c.JSON(http.StatusOK, todo)
			}
		})
```



### 删除代办

```go
		// 删除
		v1Group.DELETE("todo/:id", func(c *gin.Context) {
			// 获取前端传来的id
			id, ok := c.Params.Get("id")
			if !ok {
				c.JSON(http.StatusOK, gin.H{"error": "接收id失败"})
				return
			}

			// 删除
			err = DB.Where("id = ?", id).Delete(Todo{}).Error
			if err != nil {
				c.JSON(http.StatusOK, gin.H{"error": err.Error()})
			} else {
				c.JSON(http.StatusOK, gin.H{id: "delete成功"})
			}
		})
```



到此，所有基本功能编写完毕！



## 项目结构拆分

在企业级开发中，我们所有的代码肯定不是都放在一个程序里面的，我们需要对它进行结构拆分，让项目更加清晰，方便后续维护！

整体目录结构：

![image-20220513200927966](https://run-notes.oss-cn-beijing.aliyuncs.com/notes/02783de8f03abd1016255c11c2c07188.png)



### dao

`mysql.go`，MySQL数据库的连接、关闭，初始化DB的操作。

```go
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
```





### models

`todo.go`，models中的结构体对象和数据库的表对应，有关todo的所有操作。

```go
package models

import (
	"go-web-project-todolist/dao"
)

// Todo Model：代办清单结构体
type Todo struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Status bool   `json:"status"`
}

/*
	Todo这个Model的增删改查操作都放在这里
*/
func AddTodo(todo *Todo) (err error) {
	err = dao.DB.Create(&todo).Error
	return
}

func QueryAllTodos() (todoList []*Todo, err error) {
	if err = dao.DB.Find(&todoList).Error; err != nil {
		return nil, err
	}
	return
}

func QueryTodoById(id string) (todo *Todo, err error) {
	todo = new(Todo)
	if err = dao.DB.Debug().Where("id = ?", id).First(todo).Error; err != nil {
		return nil, err
	}
	return
}

func UpdateTodoStatusById(id string) (err error) {
	// 先查询出当前的todo事项，才知道修改成什么
	var todo Todo
	err = dao.DB.Where("id = ?", id).First(&todo).Error

	// 再根据当前todo的status去修改status
	if todo.Status {
		dao.DB.Model(&todo).Select("status").Update(map[string]interface{}{"status": false})
	} else {
		dao.DB.Model(&todo).Select("status").Update(map[string]interface{}{"status": true})
	}
	return
}

func DeleteTodoById(id string) (err error) {
	err = dao.DB.Where("id = ?", id).Delete(Todo{}).Error
	return
}
```



### controller

`controller.go`，主要与前端对接，负责接收参数和返回响应。

```go
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
```



### routers

`routers.go`，路由层，主要负责请求转发。

```go
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
```



在实际开发中，很多业务逻辑是十分复杂的，此时需要另外分出一个 `logic` 逻辑层来处理具体的逻辑服务。
