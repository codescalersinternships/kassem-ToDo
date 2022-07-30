package main

import (
	//"bytes"
	"bytes"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"github.com/gorilla/mux"
)

func MakeTempFile(t testing.TB) string {
	f, err := os.CreateTemp("", "go-sqlite-test")
	defer f.Close()
	if err != nil {
		t.Fatalf("Error making temp file: %q", err.Error())
	}
	return f.Name()
}

func TestGetALlToDo(t *testing.T) {

	t.Run("Getting all todos", func(t *testing.T) {
		file := MakeTempFile(t)
		defer os.Remove(file)

		r := mux.NewRouter()
		testApp, appErr := NewApp(file, Port, r)
		testApp.db.DB.AutoMigrate(&ToDo{})
		testApp.db.DB.Create(&ToDo{Task: "first todo"})
		testApp.db.DB.Create(&ToDo{Task: "second todo"})
		testApp.db.DB.Create(&ToDo{Task: "third todo"})
		request := httptest.NewRequest(http.MethodGet, "localhost:8080/api/todo/all", nil)
		response := httptest.NewRecorder()
		testApp.getALlToDoHandler(response, request)
		got := response.Body.String()
		want := "[{\"ID\":1,\"task\":\"first todo\",\"done\":false},{\"ID\":2,\"task\":\"second todo\",\"done\":false},{\"ID\":3,\"task\":\"third todo\",\"done\":false}]"
		if got != want || appErr != nil {
			t.Errorf("got %q want %q", got, want)
		}

	})
	t.Run("Get all to do with empty database", func(t *testing.T) {
		file := MakeTempFile(t)
		defer os.Remove(file)
		r := mux.NewRouter()
		testApp, appErr := NewApp(file, Port, r)
		testApp.db.DB.AutoMigrate(&ToDo{})
		request := httptest.NewRequest(http.MethodGet, "localhost:8080/api/todo", nil)
		response := httptest.NewRecorder()
		testApp.getALlToDoHandler(response, request)
		got := response.Body.String()
		want := "[]"
		if got != want || appErr != nil {
			t.Errorf("got %q want %q", got, want)
		}
	})
	t.Run("Get all to do with empty file", func(t *testing.T) {
		file := MakeTempFile(t)
		defer os.Remove(file)
		r := mux.NewRouter()
		testApp, appErr := NewApp(file, Port, r)
		request := httptest.NewRequest(http.MethodGet, "localhost:8080/api/todo", nil)
		response := httptest.NewRecorder()
		testApp.getALlToDoHandler(response, request)
		got := response.Result().StatusCode
		want := 500
		if got != want || appErr != nil {
			t.Errorf("got %q want %q", got, want)
		}
	})

}

func TestGetTodoHandler(t *testing.T) {
	t.Run("Get todo with existed id", func(t *testing.T) {
		file := MakeTempFile(t)
		defer os.Remove(file)
		r := mux.NewRouter()
		testApp, appErr := NewApp(file, Port, r)
		testApp.db.DB.AutoMigrate(&ToDo{})
		testApp.db.DB.Create(&ToDo{ID: 2, Task: "first todo with id 2"})
		request := httptest.NewRequest(http.MethodGet, "localhost:8080/api/todo/?taskId=2", nil)
		response := httptest.NewRecorder()
		testApp.getTodoByIdHandler(response, request)
		got := response.Body.String()
		want := "{\"ID\":2,\"task\":\"first todo with id 2\",\"done\":false}"
		if got != want || appErr != nil {
			t.Errorf("got %q want %q", got, want)
		}
	})
	t.Run("Get todo with non existed", func(t *testing.T) {
		file := MakeTempFile(t)
		defer os.Remove(file)
		r := mux.NewRouter()
		testApp, appErr := NewApp(file, Port, r)
		testApp.db.DB.AutoMigrate(&ToDo{})
		testApp.db.DB.Create(&ToDo{ID: 2, Task: "first todo with id 2"})
		testApp.db.DB.Create(&ToDo{Task: "second todo"})
		testApp.db.DB.Create(&ToDo{Task: "third todo"})
		request := httptest.NewRequest(http.MethodGet, "localhost:8080/api/todo/?taskId=24", nil)
		response := httptest.NewRecorder()
		testApp.getTodoByIdHandler(response, request)
		got := response.Body.String()
		want := "{\"response\":\"Task not found\"}\n"
		if got != want || appErr != nil {
			t.Errorf("got %q want %q", got, want)
		}
	})
	t.Run("Get todo with empty file", func(t *testing.T) {
		file := MakeTempFile(t)
		defer os.Remove(file)
		r := mux.NewRouter()
		testApp, appErr := NewApp(file, Port, r)
		request := httptest.NewRequest(http.MethodGet, "localhost:8080/api/todo/?taskId=24", nil)
		response := httptest.NewRecorder()
		testApp.getTodoByIdHandler(response, request)
		got := response.Result().StatusCode
		want := 500
		if got != want || appErr != nil {
			t.Errorf("got %q want %q", got, want)
		}
	})

}

func TestNewTaskHandler(t *testing.T) {
	t.Run("create task", func(t *testing.T) {
		file := MakeTempFile(t)
		defer os.Remove(file)
		r := mux.NewRouter()
		testApp, appErr := NewApp(file, Port, r)
		testApp.db.DB.AutoMigrate(&ToDo{})
		var todoJSON = []byte(`{"Task":"new todo"}`)
		request := httptest.NewRequest(http.MethodPost, "localhost:8080/api/todo", bytes.NewBuffer(todoJSON))
		response := httptest.NewRecorder()
		testApp.newTaskHandler(response, request)
		got := response.Body.String()
		want := "{\"ID\":1,\"task\":\"new todo\",\"done\":false}"
		if got != want || appErr != nil {
			t.Errorf("got %q want %q", got, want)
		}
	})
	t.Run("create task with empty json request", func(t *testing.T) {
		file := MakeTempFile(t)
		defer os.Remove(file)
		r := mux.NewRouter()
		testApp, appErr := NewApp(file, Port, r)
		testApp.db.DB.AutoMigrate(&ToDo{})
		var todoJSON = []byte(`{}`)
		request := httptest.NewRequest(http.MethodPost, "localhost:8080/api/todo", bytes.NewBuffer(todoJSON))
		response := httptest.NewRecorder()
		testApp.newTaskHandler(response, request)
		got := response.Result().StatusCode
		want := 400
		if got != want || appErr != nil {
			t.Errorf("got %q want %q", got, want)
		}
	})
	t.Run("decoder failure", func(t *testing.T) {
		file := MakeTempFile(t)
		defer os.Remove(file)
		r := mux.NewRouter()
		testApp, appErr := NewApp(file, Port, r)
		testApp.db.DB.AutoMigrate(&ToDo{})
		var todoJSON = []byte(``)
		request := httptest.NewRequest(http.MethodPost, "localhost:8080/api/todo", bytes.NewBuffer(todoJSON))
		response := httptest.NewRecorder()
		testApp.newTaskHandler(response, request)
		got := response.Result().StatusCode
		want := 500
		if got != want || appErr != nil {
			t.Errorf("got %v want %v", got, want)
		}
	})


}

func TestUpdateTaskHandler(t *testing.T) {
	t.Run("update task", func(t *testing.T) {
		file := MakeTempFile(t)
		defer os.Remove(file)
		r := mux.NewRouter()
		testApp, appErr := NewApp(file, Port, r)
		testApp.db.DB.AutoMigrate(&ToDo{})
		testApp.db.DB.Create(&ToDo{ID: 2, Task: "old todo"})
		var todoJSON = []byte(`{"Task":"update todo with id 2"}`)
		request := httptest.NewRequest(http.MethodPut, "localhost:8080/api/todo/?taskId=2", bytes.NewBuffer(todoJSON))
		response := httptest.NewRecorder()
		testApp.updateTaskHandler(response, request)
		got := response.Body.String()
		want := "{\"ID\":2,\"task\":\"update todo with id 2\",\"done\":false}"
		if got != want || appErr != nil {
			t.Errorf("got %q want %q", got, want)
		}
	})
	t.Run("update task with empty body request", func(t *testing.T) {
		file := MakeTempFile(t)
		defer os.Remove(file)
		r := mux.NewRouter()
		testApp, appErr := NewApp(file, Port, r)
		testApp.db.DB.AutoMigrate(&ToDo{})
		var todoJSON = []byte(`{}`)
		request := httptest.NewRequest(http.MethodPut, "localhost:8080/api/todo/?taskId=2", bytes.NewBuffer(todoJSON))
		response := httptest.NewRecorder()
		testApp.updateTaskHandler(response, request)
		got := response.Result().StatusCode
		want := 400
		if got != want || appErr != nil {
			t.Errorf("got %v want %v", got, want)
		}
	})
	t.Run("update non existed task", func(t *testing.T) {
		file := MakeTempFile(t)
		defer os.Remove(file)
		r := mux.NewRouter()
		testApp, appErr := NewApp(file, Port, r)
		testApp.db.DB.AutoMigrate(&ToDo{})
		var todoJSON = []byte(`{"Task":"new todo"}`)
		request := httptest.NewRequest(http.MethodPut, "localhost:8080/api/todo/?taskId=2233", bytes.NewBuffer(todoJSON))
		response := httptest.NewRecorder()
		testApp.updateTaskHandler(response, request)
		got := response.Result().StatusCode
		want := 400
		if got != want || appErr != nil {
			t.Errorf("got %v want %v", got, want)
		}
	})
	t.Run("decoder failure", func(t *testing.T) {
		file := MakeTempFile(t)
		defer os.Remove(file)
		r := mux.NewRouter()
		testApp, appErr := NewApp(file, Port, r)
		testApp.db.DB.AutoMigrate(&ToDo{})
		var todoJSON = []byte(``)
		request := httptest.NewRequest(http.MethodPut, "localhost:8080/api/todo", bytes.NewBuffer(todoJSON))
		response := httptest.NewRecorder()
		testApp.updateTaskHandler(response, request)
		got := response.Result().StatusCode
		want := 500
		if got != want || appErr != nil {
			t.Errorf("got %v want %v", got, want)
		}
	})
	t.Run("empty database", func(t *testing.T) {
		file := MakeTempFile(t)
		defer os.Remove(file)
		r := mux.NewRouter()
		testApp, appErr := NewApp(file, Port, r)
		var todoJSON = []byte(`{"Task":"update todo with id 2"}`)
		os.Create(file)
		request := httptest.NewRequest(http.MethodPut, "localhost:8080/api/todo/?taskId=2",bytes.NewBuffer(todoJSON))
		response := httptest.NewRecorder()
		testApp.DeleteTaskHandler(response, request)
		got := response.Result().StatusCode
		want := 500
		if got != want || appErr != nil {
			t.Errorf("got %v want %v", got, want)
		}
	})


}

func TestDeleteTaskHandler(t *testing.T) {
	t.Run("delete task", func(t *testing.T) {
		file := MakeTempFile(t)
		defer os.Remove(file)
		r := mux.NewRouter()
		testApp, appErr := NewApp(file, Port, r)
		testApp.db.DB.AutoMigrate(&ToDo{})
		testApp.db.DB.Create(&ToDo{ID: 2, Task: "old todo"})
		request := httptest.NewRequest(http.MethodPut, "localhost:8080/api/todo/?taskId=2",nil)
		response := httptest.NewRecorder()
		testApp.DeleteTaskHandler(response, request)
		got := response.Result().StatusCode
		want := 204
		if got != want || appErr != nil {
			t.Errorf("got %v want %v", got, want)
		}
	})
	t.Run("delete non existed task", func(t *testing.T) {
		file := MakeTempFile(t)
		defer os.Remove(file)
		r := mux.NewRouter()
		testApp, appErr := NewApp(file, Port, r)
		testApp.db.DB.AutoMigrate(&ToDo{})
		testApp.db.DB.Create(&ToDo{ID: 2, Task: "old todo"})
		request := httptest.NewRequest(http.MethodPut, "localhost:8080/api/todo/?taskId=2343",nil)
		response := httptest.NewRecorder()
		testApp.DeleteTaskHandler(response, request)
		got := response.Result().StatusCode
		want := 400
		if got != want || appErr != nil {
			t.Errorf("got %v want %v", got, want)
		}
	})
	t.Run("empty database", func(t *testing.T) {
		file := MakeTempFile(t)
		defer os.Remove(file)
		r := mux.NewRouter()
		testApp, appErr := NewApp(file, Port, r)
		testApp.db.DB.Create(&ToDo{ID: 2, Task: "old todo"})
		request := httptest.NewRequest(http.MethodPut, "localhost:8080/api/todo/?taskId=2",nil)
		response := httptest.NewRecorder()
		testApp.DeleteTaskHandler(response, request)
		got := response.Result().StatusCode
		want := 500
		if got != want || appErr != nil {
			t.Errorf("got %v want %v", got, want)
		}
	})

}