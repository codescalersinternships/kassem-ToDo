package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
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
		var db Database
		db.DB, _ = gorm.Open(sqlite.Open(file), &gorm.Config{})
		db.DB.AutoMigrate(&ToDo{})
		db.DB.Create(&ToDo{Task: "first todo"})
		db.DB.Create(&ToDo{Task: "second todo"})
		db.DB.Create(&ToDo{Task: "third todo"})
		request := httptest.NewRequest(http.MethodGet, "localhost:8080/todo/all", nil)
		response := httptest.NewRecorder()
		db.getALlToDo(response, request)
		got := response.Body.String()
		want := "[{\"ID\":1,\"task\":\"first todo\",\"done\":false},{\"ID\":2,\"task\":\"second todo\",\"done\":false},{\"ID\":3,\"task\":\"third todo\",\"done\":false}]\n"
		if got != want {
			t.Errorf("got %q want %q", got, want)
		}

	})
	t.Run("Get all to do with empty database", func(t *testing.T) {
		file := MakeTempFile(t)
		defer os.Remove(file)
		var db Database
		db.DB, _ = gorm.Open(sqlite.Open(file), &gorm.Config{})
		db.DB.AutoMigrate(&ToDo{})
		request := httptest.NewRequest(http.MethodGet, "localhost:8080/todo", nil)
		response := httptest.NewRecorder()
		db.getALlToDo(response, request)
		got := response.Body.String()
		want := "[]\n"
		if got != want {
			t.Errorf("got %q want %q", got, want)
		}
	})

}
func TestGetTodo(t *testing.T) {
	t.Run("Get todo with existed id", func(t *testing.T) {
		file := MakeTempFile(t)
		defer os.Remove(file)
		var db Database
		db.DB, _ = gorm.Open(sqlite.Open(file), &gorm.Config{})
		db.DB.AutoMigrate(&ToDo{})
		db.DB.Create(&ToDo{ID: 2, Task: "first todo with id 2"})
		request := httptest.NewRequest(http.MethodGet, "localhost:8080/todo/?taskId=2", nil)
		response := httptest.NewRecorder()
		db.getTodo(response, request)
		got := response.Body.String()
		want := "{\"ID\":2,\"task\":\"first todo with id 2\",\"done\":false}\n"
		if got != want {
			t.Errorf("got %q want %q", got, want)
		}
	})
	t.Run("Get todo with non existed", func(t *testing.T) {
		file := MakeTempFile(t)
		defer os.Remove(file)
		var db Database
		db.DB, _ = gorm.Open(sqlite.Open(file), &gorm.Config{})
		db.DB.AutoMigrate(&ToDo{})
		db.DB.Create(&ToDo{ID: 2, Task: "first todo with id 2"})
		db.DB.Create(&ToDo{Task: "second todo"})
		db.DB.Create(&ToDo{Task: "third todo"})
		request := httptest.NewRequest(http.MethodGet, "localhost:8080/todo/?taskId=24", nil)
		response := httptest.NewRecorder()
		db.getTodo(response, request)
		got := response.Body.String()
		want := "{\"msg\":\"Task not found\"}\n"
		if got != want {
			t.Errorf("got %q want %q", got, want)
		}
	})
}

func TestNewTask(t *testing.T) {
	t.Run("create task", func(t *testing.T) {
		file := MakeTempFile(t)
		defer os.Remove(file)
		var db Database
		db.DB, _ = gorm.Open(sqlite.Open(file), &gorm.Config{})
		db.DB.AutoMigrate(&ToDo{})
		var todoJSON = []byte(`{"Task":"new todo"}`)
		request := httptest.NewRequest(http.MethodPost, "localhost:8080/todo", bytes.NewBuffer(todoJSON))
		response := httptest.NewRecorder()
		db.newTask(response, request)
		got := response.Body.String()
		want := "{\"ID\":1,\"task\":\"new todo\",\"done\":false}\n"
		if got != want {
			t.Errorf("got %q want %q", got, want)
		}
	})
	t.Run("create task with empty body request", func(t *testing.T) {
		file := MakeTempFile(t)
		defer os.Remove(file)
		var db Database
		db.DB, _ = gorm.Open(sqlite.Open(file), &gorm.Config{})
		db.DB.AutoMigrate(&ToDo{})
		var todoJSON = []byte(`{}`)
		request := httptest.NewRequest(http.MethodPost, "localhost:8080/todo", bytes.NewBuffer(todoJSON))
		response := httptest.NewRecorder()
		db.newTask(response, request)
		got := response.Body.String()
		want := "{\"msg\":\"creation error, make sure to add task\"}\n"
		if got != want {
			t.Errorf("got %q want %q", got, want)
		}
	})
	
}
