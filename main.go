package main

import (
	"ToDoApiCheck24/viewmodel/handler"
	_ "database/sql"
	_ "fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"net/http"
)

func main() {

	//Server
	var serveMux = http.NewServeMux()
	fileServer := http.FileServer(http.Dir("/view"))
	serveMux.Handle("/", http.StripPrefix("/view", fileServer))
	serveMux.HandleFunc("/login", handler.Login)
	serveMux.HandleFunc("/register", handler.Register)
	serveMux.HandleFunc("/todo", handler.ToDo)

	log.Printf("About to listen on 8443. Go to http://127.0.0.1:8443/register\n Go to http://127.0.0.1:8443/login")
	_ = http.ListenAndServe(":8443", serveMux)

}
