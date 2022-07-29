package main

import (
	"encoding/json"
	"fmt"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func (d *database)connectDatabase(path string) error   {
	d.DB, d.err = gorm.Open(sqlite.Open(path), &gorm.Config{
		Logger :logger.Default.LogMode(logger.Silent),
	})
	
		return d.err
}

func (d *database) GetALlToDo() ([]byte, error) {
	var tasks []ToDo
	res := d.DB.Find(&tasks)
	if res.Error != nil {
		return nil, res.Error
	}
	data, err := json.Marshal(tasks)
	return data, err
}

func (d *database) GetTodoById(id string) ([]byte, error) {
	var res ToDo
	var data []byte
	var err error

	query := d.DB.First(&res, id)
	if query.Error != nil {
 		data =nil
		err = query.Error

	} else {
		data, err = json.Marshal(res)
	}
	return data, err
}
func (d *database) AddTodo(task *ToDo) ([]byte, error) {
	var data []byte
	var err error
	query := d.DB.Create(&task)
	
	if query.Error != nil {
		data =nil
		err = query.Error
	} else {
		data,err =json.Marshal(&task)
	}
	return data, err
}
func (d *database) UpdateTodo(NewTask *ToDo, id string) ([]byte, error) {
	var res ToDo
	var data []byte
	var err error

	query := d.DB.First(&res, id)
	if query.Error != nil {
 		data =nil
		err = query.Error
		return data, err
	} 
	res.Task = NewTask.Task
	res.Done = NewTask.Done
 	updateQuery := d.DB.Save(&res)
	if updateQuery.Error != nil {
		data=nil
		err = updateQuery.Error
	}else {
		data, err = json.Marshal(&res)
	}	
	return data, err
}
func (d *database) DeleteTask(id string) ([]byte, error) {
	var data []byte
	var err error
	query:=d.DB.Delete(&ToDo{}, id)
	if query.Error != nil {
		data = nil
		err = query.Error
	}else if query.RowsAffected == 0{
		err = gorm.ErrRecordNotFound
		data = nil
	}else {
		data, err = json.Marshal(fmt.Sprintf("{ response: task with id %v deleted successfully }",id))
	}
	return data, err
}