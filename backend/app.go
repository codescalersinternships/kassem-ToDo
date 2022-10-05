package main

import (
	"encoding/json"

	"fmt"
	"net/http"
	
	"strconv"

	"gorm.io/gorm"
)

func (a *App) home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:8080")
	w.Header().Set("Access-Control-Allow-Headers", "authentication, content-type")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode("ToDo home")
	w.WriteHeader(http.StatusOK)
}

func (a *App) getALlToDoHandler(w http.ResponseWriter, r *http.Request) {
	//w.Header().Set("Access-Control-Allow-Origin", "http://localhost:8080")
	w.Header().Set("Content-Type", "application/json")

	res, err := a.db.GetALlToDo()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode("Error :" + err.Error())

		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(res)

}

func (a *App) getTodoByIdHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:8080")
	w.Header().Set("Content-Type", "application/json")

	id := r.URL.Query().Get("taskId")
	res, err := a.db.GetTodoById(id)
	if err == gorm.ErrRecordNotFound {
		errResp := Response{Response: "Task not found"}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errResp)
	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(&err)
	} else {
		w.WriteHeader(http.StatusOK)
		w.Write(res)
	}
}

// add new task to todo database
func (a *App) newTaskHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:80")
	w.Header().Set("Access-Control-Allow-Methods", "*")
	w.Header().Set("Content-Type", "application/json")
	var task ToDo
	//get body data
	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	resp, respErr := a.db.AddTodo(&task)
	if respErr != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(respErr)
	} else {
		w.WriteHeader(http.StatusOK)
		w.Write(resp)
	}

}

func (a *App) updateTaskHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")                                                            // 允许访问所有域，可以换成具体url，注意仅具体url才能带cookie信息
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, Token") //header的类型
	w.Header().Add("Access-Control-Allow-Credentials", "true")                                                    //设置为true，允许ajax异步请求带cookie信息
	w.Header().Add("Access-Control-Allow-Methods", "PUT, GET, OPTIONS, POST, DELETE")                             //允许请求方法
	w.Header().Set("content-type", "application/json;charset=UTF-8")

	id := r.URL.Query().Get("taskId")
	var tmp ToDo
	tmp.ID, _ = strconv.Atoi(id)
	err := json.NewDecoder(r.Body).Decode(&tmp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err.Error())
		return
	}
	if tmp.Task == "" {
		errResp := Response{Response: "Make sure to add task"}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errResp)
		return
	}
	res, err := a.db.UpdateTodo(&tmp, id)
	if err == gorm.ErrRecordNotFound {
		errResp := Response{Response: "Task not found"}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errResp)
	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(&err)
	} else {
		w.WriteHeader(http.StatusOK)
		w.Write(res)
	}

}

func (a *App) DeleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id := r.URL.Query().Get("taskId")
	res, err := a.db.DeleteTask(id)

	if err == gorm.ErrRecordNotFound {
		errResp := Response{Response: "Task not found"}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errResp)
	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err.Error())
		return
	} else {
		w.WriteHeader(http.StatusNoContent)
		w.Write(res)
	}
}

func NewApp(DB_FILE, Port string, router http.Handler) (App, error) {
	a := App{}
	a.db = &database{}
	a.server = &http.Server{}
	err := a.db.connectDatabase(DB_FILE)
	fmt.Printf("start sever on port: %s\n", Port)
	
	a.db.DB.AutoMigrate(&ToDo{})
	a.server.Addr = Port
	a.server.Handler = router
	return a, err
}
