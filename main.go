package main 

import (
	"github.com/gorilla/mux"
	"net/http"
	"log"
	"time"
)
var DB_FILE = "todo.db"
var Port = ":8080"

func logTimeMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%v: %v\n%v", r.Method, r.RequestURI,time.Now().Format(time.RFC850))
		handler.ServeHTTP(w, r)
	})
}

func main() {
	r := mux.NewRouter()
	app,err := NewApp(DB_FILE, Port,r)
	if err != nil {
		panic("couldn't connect to database path" + DB_FILE + " due to error: " + err.Error() )
	}
	r.Use(logTimeMiddleware)
	r.HandleFunc("/home",app.home).Methods("GET")
	r.HandleFunc("/api/todo/all", app.getALlToDoHandler).Methods("GET")
	r.HandleFunc("/api/todo/", app.getTodoByIdHandler).Methods("GET")
	r.HandleFunc("/api/todo", app.newTaskHandler).Methods("POST")
	r.HandleFunc("/api/todo/", app.updateTaskHandler).Methods("PUT")
	r.HandleFunc("/api/todo/", app.DeleteTaskHandler).Methods("DELETE")

	app.server = &http.Server{Addr: Port, Handler: r}
	app.server.ListenAndServe();
}