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
	serveMux.HandleFunc("/delete", handler.Delete)
	serveMux.HandleFunc("/update", handler.Update)
	serveMux.HandleFunc("/logout", handler.Logout)
	serveMux.HandleFunc("/addUser", handler.AddUser)
	log.Printf("About to listen on 8443. Go to http://127.0.0.1:8443/register\n Go to http://127.0.0.1:8443/login")
	//err := http.ListenAndServeTLS(":8443", "OpenPGP_0xDC5E2C7D7B83A73E.asc", "OpenPGP_signature", serveMux)
	err := http.ListenAndServe(":8443", serveMux)
	if err != nil {
		log.Fatal(err)
	}

}
