package handler

import (
	"ToDoApiCheck24/model"
	"database/sql"
	"fmt"
	"html/template"
	_ "html/template"
	"net/http"
	"strconv"
)

func Login(responseWriter http.ResponseWriter, request *http.Request) {
	if request.Method == "POST" {

		var users []model.User
		users = getAllUSersFromDB()
		for _, user := range users {
			if user.Username == request.FormValue("username") && user.Password == request.FormValue("password") {
				//Login Success!
				setCookie(responseWriter, "username", user.Username)
				setCookie(responseWriter, "id", strconv.Itoa(user.IdUsers))
				http.Redirect(responseWriter, request, "/todo", 301)
				return
			} else {
				//Wrong Username or Password
				parseAndExecuteWebsite("view/html/login.html", responseWriter, "Wrong Username or Password!")
				return

			}
		}
	}
	if request.Method == "GET" {
		parseAndExecuteWebsite("view/html/login.html", responseWriter, nil)
	}
}

func Register(responseWriter http.ResponseWriter, request *http.Request) {
	if request.Method == "POST" {

		var users []model.User
		users = getAllUSersFromDB()
		for _, user := range users {
			if user.Username == request.FormValue("username") {
				fmt.Println("Username Taken!")
				parseAndExecuteWebsite("view/html/register.html", responseWriter, "Username already Taken!")
				return
			} else {
				db, err := sql.Open("mysql", "root:admin@tcp(127.0.0.1:3306)/todoapi")
				if err != nil {
					panic(err.Error())
				}
				_, err = db.Exec("INSERT INTO `todoapi`.`users` (`Username`, `Password`) VALUES ('" + request.FormValue("username") + "', '" + request.FormValue("password") + "');")
				if err != nil {
					return
				}
				defer closeDB(db)
				fmt.Println("Success!")
				http.Redirect(responseWriter, request, "/login", 301)

			}
		}

	}

	if request.Method == "GET" {
		parseAndExecuteWebsite("view/html/register.html", responseWriter, nil)
	}

}

func ToDo(responseWriter http.ResponseWriter, request *http.Request) {
	if request.Method == "POST" {
		db, err := sql.Open("mysql", "root:admin@tcp(127.0.0.1:3306)/todoapi")
		if err != nil {
			panic(err.Error())
		}
		todosName := request.FormValue("todo")
		_, err = db.Exec("INSERT INTO `todoapi`.`todos` (`TodosName`, `TodosDone`) VALUES ('" + todosName + "', '0');")
		defer closeDB(db)

		var todos []model.ToDo
		todos = getAllTodosFromDB()
		userid := informationsFromCookies("id", request)
		db, err = sql.Open("mysql", "root:admin@tcp(127.0.0.1:3306)/todoapi")
		idTodos := strconv.Itoa(todos[len(todos)-1].IdTodos)
		queryInsert := "INSERT INTO `todoapi`.`todoowners` (`IdOfOwner`, `IdOfTodo`) VALUES ('" + userid + "', '" + idTodos + "');"
		_, err = db.Exec(queryInsert)
		if err != nil {
			return
		}
		defer closeDB(db)
		//http.Redirect(responseWriter, request, "/todo", 301)
		parseAndExecuteWebsite("view/html/todo.html", responseWriter, getAllTodosFromUsers(request))
		return
	}

	if request.Method == "GET" {

		parseAndExecuteWebsite("view/html/todo.html", responseWriter, getAllTodosFromUsers(request))
	}

}

func getAllTodosFromUsers(request *http.Request) []model.ToDo {
	var toDos []model.ToDo
	var toDoOwner []model.TodoOwner
	userid := informationsFromCookies("id", request)
	toDos = getAllTodosFromDB()
	toDoOwner = getAllTodoOwnerFromDB()
	var todosNumbers []int
	for _, do := range toDoOwner {
		if strconv.Itoa(do.IdOfOwner) == userid {
			todosNumbers = append(todosNumbers, do.IdOfTodo)
		}
	}
	var todosResult []model.ToDo
	for _, number := range todosNumbers {
		for _, do := range toDos {
			if do.IdTodos == number {
				todosResult = append(todosResult, do)
			}
		}
	}
	return todosResult
}
func informationsFromCookies(value string, request *http.Request) string {
	for _, cookie := range request.Cookies() {
		if cookie.Name == "username" && value == "username" {
			return cookie.Value
		}
		if cookie.Name == "id" && value == "id" {
			return cookie.Value
		}
	}
	return "NO INFORMATION"
}

func getAllUSersFromDB() []model.User {
	db, err := sql.Open("mysql", "root:admin@tcp(127.0.0.1:3306)/todoapi")
	if err != nil {
		panic(err.Error())
	}
	result, err := db.Query("SELECT IdUsers, Username, Password FROM users")
	if err != nil {
		panic(err.Error())
	}
	var users []model.User
	for result.Next() {
		var user model.User
		err = result.Scan(&user.IdUsers, &user.Username, &user.Password)
		if err != nil {
			panic(err.Error())
		}
		users = append(users, user)
	}
	defer closeDB(db)
	return users
}
func getAllTodosFromDB() []model.ToDo {
	db, err := sql.Open("mysql", "root:admin@tcp(127.0.0.1:3306)/todoapi")
	if err != nil {
		panic(err.Error())
	}
	result, err := db.Query("SELECT IdTodos,TodosName, TodosDone FROM todos")
	if err != nil {
		panic(err.Error())
	}
	var todos []model.ToDo
	for result.Next() {
		var toDo model.ToDo
		err = result.Scan(&toDo.IdTodos, &toDo.TodosName, &toDo.TodosDone)
		if err != nil {
			panic(err.Error())
		}
		todos = append(todos, toDo)
	}
	print(result)
	defer closeDB(db)
	return todos
}
func getAllTodoOwnerFromDB() []model.TodoOwner {
	db, err := sql.Open("mysql", "root:admin@tcp(127.0.0.1:3306)/todoapi")
	if err != nil {
		panic(err.Error())
	}
	result, err := db.Query("SELECT IdTodoOwners,IdOfOwner,IdOfTodo FROM todoowners")
	if err != nil {
		panic(err.Error())
	}
	var todoOwners []model.TodoOwner
	for result.Next() {
		var todoOwner model.TodoOwner
		err = result.Scan(&todoOwner.IdTodoOwner, &todoOwner.IdOfOwner, &todoOwner.IdOfTodo)
		if err != nil {
			panic(err.Error())
		}
		todoOwners = append(todoOwners, todoOwner)
	}
	defer closeDB(db)
	return todoOwners
}
func closeDB(db *sql.DB) {
	err := db.Close()
	if err != nil {

	}
}

func parseAndExecuteWebsite(filename string, responseWriter http.ResponseWriter, data interface{}) {
	t, err := template.ParseFiles(filename)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = t.Execute(responseWriter, data)
	if err != nil {
		fmt.Println(err)
		return
	}
}
func setCookie(responseWriter http.ResponseWriter, name string, value string) {

	cookieToStore := http.Cookie{Name: name, Value: value}
	http.SetCookie(responseWriter, &cookieToStore)
}
