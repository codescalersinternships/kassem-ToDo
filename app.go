package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

//gorm db

// "root:password@tcp(127.0.0.1:3306)/Todo?charset=utf8mb4&parseTime=True&loc=Local"
var dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
	os.Getenv("MYSQL_USER"), os.Getenv("MYSQL_PASSWORD"), os.Getenv("MYSQL_HOST"), os.Getenv("MYSQL_PORT"), os.Getenv("MYSQL_DATABASE"))

var db, errorDB = gorm.Open(mysql.Open(dsn), &gorm.Config{})

type Todo struct {
	ID     int    `gorm:"autoIncrement" json:"id"`
	Task   string `json:"task"`
	Status bool   `json:"done"`
}

func home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode("Todo home")
	w.WriteHeader(http.StatusOK)
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
	id, _ := strconv.Atoi(params["taskId"])
	var tasks = Todo{ID: id}
	db.Find(&tasks, id)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&tasks)

}
func newTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var task Todo
	//get body data
	_ = json.NewDecoder(r.Body).Decode(&task)
	db.Create(&task)
	if errorDB == nil {
		json.NewEncoder(w).Encode(&task)
		w.WriteHeader(http.StatusOK)
	} else {
		log.Fatalln(db.Error)
	}

}

func updateTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	//convert params id to int
	id, _ := strconv.Atoi(params["taskId"])
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
	id, _ := strconv.Atoi(params["taskId"])
	var task = Todo{ID: id}
	db.Delete(&task)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("200 - Task deleted"))
}

func main() {
	fmt.Println("editor:", os.Getenv("EDITOR"))
	fmt.Println(dsn)
	//init route
	mux := mux.NewRouter()
	// if errorDB != nil {
	// 	log.Fatal(errorDB)
	// }
	// create table if not exists
	if !(db.Migrator().HasTable(&Todo{})) {
		log.Println("table { todos } created")
		db.Migrator().CreateTable(&Todo{})
	}

	// route endpoint handler
	mux.HandleFunc("/", home)
	mux.HandleFunc("/api/todos", todos).Methods("GET")
	mux.HandleFunc("/api/todo/{taskId}", todo).Methods("GET")
	mux.HandleFunc("/api/todo", newTask).Methods("POST")
	mux.HandleFunc("/api/todo/{taskId}", updateTask).Methods("PUT")
	mux.HandleFunc("/api/todo/{taskId}", remove).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":8081", mux))

}
