package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

//gorm db
type Database struct {
	DB *gorm.DB
}

var DB_FILE = os.Getenv("DB_FILE")

type ToDo struct {
	ID   int    `gorm:"primaryKey"`
	Task string `gorm:"not null;default:null" json:"task"`
	Done bool   `json:"done"`
}

type message struct {
	MSG string `json:"msg"`
}

func home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode("ToDo home")
	w.WriteHeader(http.StatusOK)
}

//get all tasks in todo list
func (db *Database) getALlToDo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var tasks []ToDo
	res := db.DB.Find(&tasks)
	if res.Error != nil {
		json.NewEncoder(w).Encode("Error :" + res.Error.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
	w.WriteHeader(http.StatusOK)

}

//get task from database by id
func (db *Database) getTodo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	var res ToDo
	id := params["taskId"]
	findErr := db.DB.Find(&res, id).Error
	if findErr == nil {
		if res.ID != 0 {
			json.NewEncoder(w).Encode(&res)
			w.WriteHeader(http.StatusOK)
			return
		} else {
			dd := message{MSG: "Task not found"}

			json.NewEncoder(w).Encode(dd)
			w.WriteHeader(http.StatusOK)
			return
		}

	} else {
		log.Println(error(findErr))
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(findErr)
		return
	}
}

// add new task to todo database
func (db *Database) newTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var task ToDo
	//get body data
	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		log.Fatalln(err.Error())
		json.NewEncoder(w).Encode(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	creationErr := db.DB.Create(&task).Error
	if creationErr != nil {
		log.Println(error(creationErr))
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(creationErr)
		return
	} else {
		json.NewEncoder(w).Encode(&task)
		w.WriteHeader(http.StatusOK)
		return

	}

}

func (db *Database) updateTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	//storing the updated info in tmp
	var tmp ToDo

	err := json.NewDecoder(r.Body).Decode(&tmp)
	if err != nil {
		log.Fatalln(err.Error())
		json.NewEncoder(w).Encode(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	var task ToDo
	findErr := db.DB.Find(&task, params["taskId"]).Error

	if findErr != nil {
		log.Println(error(findErr))
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(findErr)
		return
	}
	id, _ := strconv.Atoi(params["taskId"])
	if id != task.ID {
		dd := message{MSG: "Task not found"}
		json.NewEncoder(w).Encode(dd)
		w.WriteHeader(http.StatusOK)
		return

	}
	task.Task = tmp.Task
	task.Done = tmp.Done
	updateErr := db.DB.Save(&task).Error

	if updateErr != nil {
		log.Println(error(updateErr))
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(updateErr)
		return
	} else {
		json.NewEncoder(w).Encode(&task)
		w.WriteHeader(http.StatusOK)
		return
	}
}

func (db *Database) removeTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	//get body data
	params := mux.Vars(r)
	id := params["taskId"]
	DeleteErr := db.DB.Delete(&ToDo{}, id)
	if DeleteErr.Error != nil {
		log.Println(error(DeleteErr.Error))
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(DeleteErr.Error)
		return
	} else if DeleteErr.RowsAffected != 1 {
		//json.NewEncoder(w).Encode(&task)
		dd := message{MSG: "Task not found"}
		json.NewEncoder(w).Encode(dd)
		w.WriteHeader(http.StatusOK)
		return
	} else {
		dd := message{MSG: "Task with id: " + id + ", deleted successfully"}
		json.NewEncoder(w).Encode(dd)
		w.WriteHeader(http.StatusOK)
		return
	}

}

func main() {
	db := Database{}
	var err error
	db.DB, err = gorm.Open(sqlite.Open(DB_FILE), &gorm.Config{})
	log.Println(DB_FILE)
	if err != nil {
		panic("couldn't connect")
	}
	db.DB.AutoMigrate(&ToDo{})
	//init route
	mux := mux.NewRouter()

	// route endpoint handler
	mux.HandleFunc("/", home)
	mux.HandleFunc("/api/todos", db.getALlToDo).Methods("GET")
	mux.HandleFunc("/api/todo/{taskId}", db.getTodo).Methods("GET")
	mux.HandleFunc("/api/todo", db.newTask).Methods("POST")
	mux.HandleFunc("/api/todo/{taskId}", db.updateTask).Methods("PUT")
	mux.HandleFunc("/api/todo/{taskId}", db.removeTask).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":8080", mux))

}
