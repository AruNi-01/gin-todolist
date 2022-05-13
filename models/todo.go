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
