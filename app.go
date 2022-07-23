package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Todo struct {
	Id    string `json:"id"`
	Task  string `json:"task"`
	State bool   `json:"done"`
}

// get all tasks in todo list
func todos(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(list)

}
func todo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	for _, item := range list {
		if item.Id == params["id"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
	json.NewEncoder(w).Encode(&Todo{})

}
func newTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var task Todo
	_ = json.NewDecoder(r.Body).Decode(&task)
	task.Id = strconv.Itoa(rand.Intn(100000)) //mock id
	list = append(list, task)
	json.NewEncoder(w).Encode(task)

}
func updateTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	var tmp Todo
	_ = json.NewDecoder(r.Body).Decode(&tmp)
	for index, task := range list {
		if task.Id == params["id"] {
			list = append(list[:index], list[index+1:]...)
			tmp.Id = params["id"]
			list = append(list, tmp)
			json.NewEncoder(w).Encode(tmp)
			break
		}
	}

}
func remove(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for index, item := range list {
		if item.Id == params["id"] {
			list = append(list[:index], list[index+1:]...)
			break
		}
	}
	json.NewEncoder(w).Encode(list)
}

var list []Todo

func main() {
	//init route
	mux := mux.NewRouter()

	// mock data
	list = append(list, Todo{Id: "1", Task: "test", State: false}, Todo{Id: "3", Task: "test@2", State: false})

	// route endpoint handler
	mux.HandleFunc("/api/todos", todos).Methods("GET")
	mux.HandleFunc("/api/todo/{id}", todo).Methods("GET")
	mux.HandleFunc("/api/todo", newTask).Methods("POST")
	mux.HandleFunc("/api/todo/{id}", updateTask).Methods("PUT")
	mux.HandleFunc("/api/todo/{id}", remove).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":8080", mux))

}
