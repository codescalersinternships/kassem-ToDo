package main

import (
	"encoding/json"
	"log"
	"net/http"

	//"strconv"
	// "github.com/gorilla/mux"
	//	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

//gorm db

func (a *App) home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode("ToDo home")
	w.WriteHeader(http.StatusOK)
}

func (a *App) getALlToDoHandler(w http.ResponseWriter, r *http.Request) {
	// 	w.Header().Set("Content-Type", "application/json")

	res, err := a.db.GetALlToDo()
	if err != nil {
		json.NewEncoder(w).Encode("Error :" + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)

}

func (a *App) getTodoByIdHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := r.URL.Query().Get("taskId")
	res, err := a.db.GetTodoById(id)
	if err == gorm.ErrRecordNotFound {
		errResp := Response{Response: "Task not found"}
		w.WriteHeader(http.StatusNotFound)
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
	w.Header().Set("Content-Type", "application/json")
	id := r.URL.Query().Get("taskId")
	var tmp ToDo
	err := json.NewDecoder(r.Body).Decode(&tmp)
	if err != nil {
		log.Fatalln(err.Error())
		json.NewEncoder(w).Encode(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	res, err := a.db.UpdateTodo(&tmp, id)
	if err == gorm.ErrRecordNotFound {
		errResp := Response{Response: "Task not found"}
		w.WriteHeader(http.StatusNotFound)
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

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err.Error())
		return
	} else if err == gorm.ErrRecordNotFound {
		errResp := Response{Response: "Task not found"}
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(errResp)
	} else {
		w.Write(res)
		w.WriteHeader(http.StatusNoContent)

	}
}

func logTimeMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("fefefe")
		handler.ServeHTTP(w, r)
	})
}

func NewApp(DB_FILE, Port string, router http.Handler) (App, error) {
	a := App{}
	a.db = &database{}
	a.server = &http.Server{}
	err := a.db.connectDatabase(DB_FILE)
	if err != nil {
		return a, err
	}
	a.server.Addr = Port
	a.server.Handler = router
	return a, nil
}
