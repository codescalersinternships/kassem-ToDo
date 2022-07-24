package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"github.com/gorilla/mux"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)
//gorm db
var dsn = "root:password@tcp(127.0.0.1:3306)/Todo?charset=utf8mb4&parseTime=True&loc=Local"
var db, _ = gorm.Open(mysql.Open(dsn), &gorm.Config{})

type Todo struct {
	ID     int    `gorm:"autoIncrement" json:"id"`
	Task   string `json:"task"`
	Status bool   `json:"done"`
}

//get all tasks in todo list
func todos(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	var tasks []Todo
	db.Find(&tasks)

	json.NewEncoder(w).Encode(tasks)
	w.WriteHeader(http.StatusOK)

}
func todo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	//convert params id to int
	id, _ := strconv.Atoi(params["id"])
	var tasks = Todo{ID: id}
	db.Find(&tasks, id)
	json.NewEncoder(w).Encode(&tasks)
	w.WriteHeader(http.StatusOK)

}
func newTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var task Todo
	//get body data
	_ = json.NewDecoder(r.Body).Decode(&task)
	db.Create(&task)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&task)

}

func updateTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	//convert params id to int
	id, _ := strconv.Atoi(params["id"])
	var task = Todo{ID: id}
	db.Find(&task)
	_ = json.NewDecoder(r.Body).Decode(&task)
	db.Save(&task)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&task)
}

func remove(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	//convert params id to int
	id, _ := strconv.Atoi(params["id"])
	var task = Todo{ID: id}
	db.Delete(&task)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("200 - Task deleted"))
}

func main() {
	//init route
	mux := mux.NewRouter()

	// create table if not exists
	if !(db.Migrator().HasTable(&Todo{})) {
		log.Println("table { todos } created")
		db.Migrator().CreateTable(&Todo{})
	}

	// route endpoint handler
	mux.HandleFunc("/api/todos", todos).Methods("GET")
	mux.HandleFunc("/api/todo/{id}", todo).Methods("GET")
	mux.HandleFunc("/api/todo", newTask).Methods("POST")
	mux.HandleFunc("/api/todo/{id}", updateTask).Methods("PUT")
	mux.HandleFunc("/api/todo/{id}", remove).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":8080", mux))

}
