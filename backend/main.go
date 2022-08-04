package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/joho/godotenv"
)




func logTimeMiddleware(handler http.Handler) http.Handler {
	
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%v: %v\n%v", r.Method, r.RequestURI, time.Now().Format(time.RFC850))
		handler.ServeHTTP(w, r)
	})
}

func main() {
	err := godotenv.Load()

   
	DB_FILE := os.Getenv("DB_FILE")
	if DB_FILE == "" {
		DB_FILE ="./todo.db"
	}
	Port := os.Getenv("Port")
	if Port == "" {
		Port = ":5000"
	}
	
	r := mux.NewRouter()
	c := cors.AllowAll()
	handler := c.Handler(r)
	app, err := NewApp(DB_FILE, Port, handler)
	if err != nil {
		panic("couldn't connect to database path" + DB_FILE + " due to error: " + err.Error())
	} else {
		fmt.Printf("Success, listing on port %v\n", Port)
	}
	// r.Use(corsMiddleware)
	r.Use(logTimeMiddleware)
	r.HandleFunc("/home", app.home).Methods("GET")
	r.HandleFunc("/api/todo/all", app.getALlToDoHandler).Methods("GET")
	r.HandleFunc("/api/todo/", app.getTodoByIdHandler).Methods("GET")
	r.HandleFunc("/api/todo", app.newTaskHandler).Methods("POST")
	r.HandleFunc("/api/todo/", app.updateTaskHandler).Methods("PUT")
	r.HandleFunc("/api/todo/", app.DeleteTaskHandler).Methods("DELETE")
	sh := http.StripPrefix("/swaggerui/", http.FileServer(http.Dir("./swaggerui/dist/")))
	r.PathPrefix("/swaggerui/").Handler(sh)

	app.server = &http.Server{Addr: Port, Handler: handler}

	if err := app.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("listen: %s\n", err)
	}
}
