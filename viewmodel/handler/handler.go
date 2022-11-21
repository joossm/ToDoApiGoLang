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
		users, _, _ = getFromDatabase("users")
		for _, user := range users {
			if user.Username == request.FormValue("username") &&
				user.Password == request.FormValue("password") {
				//Login Success!
				setCookie(responseWriter, "username", user.Username)
				setCookie(responseWriter, "id", strconv.Itoa(user.IdUsers))
				http.Redirect(responseWriter, request, "/todo", 301)
				return
			}
		}
		//Wrong Username or Password
		parseAndExecuteWebsite("view/html/login.html", responseWriter,
			"Wrong Username or Password!")
		return

	}
	if request.Method == "GET" {
		parseAndExecuteWebsite("view/html/login.html", responseWriter, nil)
	}
}

func Register(responseWriter http.ResponseWriter, request *http.Request) {
	if request.Method == "POST" {
		var users []model.User
		users, _, _ = getFromDatabase("users")
		for _, user := range users {
			if user.Username == request.FormValue("username") {
				fmt.Println("Username Taken!")
				parseAndExecuteWebsite("view/html/register.html", responseWriter,
					"Username already Taken!")
				return
			}
		}
		db := openDB()
		//_, err = db.Exec("INSERT INTO `todoapi`.`users` (`Username`, `Password`) VALUES ('" + request.FormValue("username") + "', '" + request.FormValue("password") + "');")
		_, err := db.Exec("INSERT INTO `todoapi`.`users` (`Username`, `Password`) VALUES (?, ?);", request.FormValue("username"), request.FormValue("password"))
		errorHandler(err)
		defer closeDB(db)
		fmt.Println("Success!")
		http.Redirect(responseWriter, request, "/login", 301)
	}
	if request.Method == "GET" {
		parseAndExecuteWebsite("view/html/register.html", responseWriter, nil)
	}
}
func AddUser(responseWriter http.ResponseWriter, request *http.Request) {
	// get id from url
	id := request.URL.Query().Get("id")
	// get id from todoapi.users
	db := openDB()
	result, err := db.Query("SELECT IdUsers, Username FROM todoapi.users WHERE Username = ?", request.FormValue("username"))
	errorHandler(err)
	defer closeDB(db)
	var user model.User
	for result.Next() {
		err := result.Scan(&user.IdUsers, &user.Username)
		errorHandler(err)
	}
	// add userid and todoid to todoapi.todoowners
	db = openDB()
	_, err = db.Exec("INSERT INTO `todoapi`.`todoowners` (`IdOfOwner`, `IdOfTodo`) VALUES (?,?);", user.IdUsers, id)
	errorHandler(err)
	defer closeDB(db)
	// redirect to /todo
	parseAndExecuteWebsite("view/html/todo.html", responseWriter, getAllTodosFromUsers(request))
}
func ToDo(responseWriter http.ResponseWriter, request *http.Request) {
	userid := informationsFromCookies("id", request)
	if userid == "NO INFORMATION" || userid == "" {
		http.Redirect(responseWriter, request, "/login", 301)
		return
	} else {
		if request.Method == "POST" {
			db := openDB()
			todosName := request.FormValue("todo")
			todosText := request.FormValue("text")
			_, err := db.Exec("INSERT INTO `todoapi`.`todos` (`TodosName`, `TodosDone`, `TodosText`) VALUES (?, '0', ?);", todosName, todosText)
			errorHandler(err)
			defer closeDB(db)
			var todos []model.ToDo
			_, todos, _ = getFromDatabase("todos")
			userid := informationsFromCookies("id", request)
			db = openDB()
			idTodos := strconv.Itoa(todos[len(todos)-1].IdTodos)
			_, err = db.Exec("INSERT INTO `todoapi`.`todoowners` (`IdOfOwner`, `IdOfTodo`) VALUES (?,?);", userid, idTodos)
			errorHandler(err)
			defer closeDB(db)
			//http.Redirect(responseWriter, request, "/todo", 301)
			parseAndExecuteWebsite("view/html/todo.html", responseWriter, getAllTodosFromUsers(request))
			return
		}
		if request.Method == "GET" {
			parseAndExecuteWebsite("view/html/todo.html", responseWriter, getAllTodosFromUsers(request))
		}

	}
}

func Delete(responseWriter http.ResponseWriter, request *http.Request) {
	userid := informationsFromCookies("id", request)
	if userid == "NO INFORMATION" || userid == "" {
		http.Redirect(responseWriter, request, "/login", 301)
		return
	} else {
		// get id from url
		id := request.URL.Query().Get("id")
		// delete todo from database
		db := openDB()
		_, err := db.Exec("DELETE FROM `todoapi`.`todos` WHERE (`IdTodos` = ?);", id)
		_, err = db.Exec("DELETE FROM `todoapi`.`todoowners` WHERE (`IdOfTodo` = ?);", id)
		errorHandler(err)
		parseAndExecuteWebsite("view/html/todo.html", responseWriter, getAllTodosFromUsers(request))
		return
	}
}
func Update(responseWriter http.ResponseWriter, request *http.Request) {
	userid := informationsFromCookies("id", request)
	if userid == "NO INFORMATION" || userid == "" {
		http.Redirect(responseWriter, request, "/login", 301)
		return
	} else {
		// get id from url
		id := request.URL.Query().Get("id")
		// update todo from database
		var todos []model.ToDo
		_, todos, _ = getFromDatabase("todos")
		for _, todo := range todos {
			if strconv.Itoa(todo.IdTodos) == id {
				if todo.TodosDone == 0 {
					db := openDB()
					_, err := db.Exec("UPDATE `todoapi`.`todos` SET `TodosDone` = '1' WHERE (`IdTodos` = ?);", id)
					errorHandler(err)
					parseAndExecuteWebsite("view/html/todo.html", responseWriter, getAllTodosFromUsers(request))
					return
				} else {
					db := openDB()
					_, err := db.Exec("UPDATE `todoapi`.`todos` SET `TodosDone` = '0' WHERE (`IdTodos` = ?);", id)
					errorHandler(err)
					parseAndExecuteWebsite("view/html/todo.html", responseWriter, getAllTodosFromUsers(request))
					return
				}
			}
		}
		db := openDB()
		_, err := db.Exec("UPDATE FROM `todoapi`.`todos` WHERE (`IdTodos` = ?);", id)
		errorHandler(err)
		parseAndExecuteWebsite("view/html/todo.html", responseWriter, getAllTodosFromUsers(request))
		return
	}
}
func Logout(responseWriter http.ResponseWriter, request *http.Request) {
	// delete cookie
	setCookie(responseWriter, "username", "")
	setCookie(responseWriter, "id", "")
	http.Redirect(responseWriter, request, "/login", 301)
	return

}

func getAllTodosFromUsers(request *http.Request) []model.ToDo {
	userid := informationsFromCookies("id", request)
	db := openDB()
	result, err := db.Query("SELECT todos.IdTodos, todos.TodosName, todos.TodosDone, todos.TodosText FROM todoapi.todoowners, todoapi.todos WHERE todoowners.IdOfTodo = todos.IdTodos AND todoowners.IdOfOwner = ?;", userid)
	errorHandler(err)
	var toDos []model.ToDo
	for result.Next() {
		var toDo model.ToDo
		err = result.Scan(&toDo.IdTodos, &toDo.TodosName, &toDo.TodosDone, &toDo.TodosText)
		errorHandler(err)
		toDos = append(toDos, toDo)
	}
	return toDos
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
func getFromDatabase(want string) ([]model.User, []model.ToDo, []model.TodoOwner) {
	db := openDB()
	if want == "users" {
		result, err := db.Query("SELECT IdUsers, Username, Password FROM users")
		errorHandler(err)
		var users []model.User
		for result.Next() {
			var user model.User
			err = result.Scan(&user.IdUsers, &user.Username, &user.Password)

			users = append(users, user)
		}
		closeDB(db)
		return users, nil, nil
	}
	if want == "todos" {
		result, err := db.Query("SELECT IdTodos,TodosName, TodosDone, TodosText FROM todos")
		errorHandler(err)
		var todos []model.ToDo
		for result.Next() {
			var toDo model.ToDo
			err = result.Scan(&toDo.IdTodos, &toDo.TodosName, &toDo.TodosDone, &toDo.TodosText)
			errorHandler(err)
			todos = append(todos, toDo)
		}
		closeDB(db)
		return nil, todos, nil
	}
	if want == "todoowners" {
		result, err := db.Query("SELECT IdTodoOwners,IdOfOwner,IdOfTodo FROM todoowners")
		errorHandler(err)
		var todoOwners []model.TodoOwner
		for result.Next() {
			var todoOwner model.TodoOwner
			err = result.Scan(&todoOwner.IdTodoOwner, &todoOwner.IdOfOwner, &todoOwner.IdOfTodo)
			errorHandler(err)
			todoOwners = append(todoOwners, todoOwner)
		}
		closeDB(db)
		return nil, nil, todoOwners
	}
	defer closeDB(db)
	return nil, nil, nil
}

func parseAndExecuteWebsite(filename string, responseWriter http.ResponseWriter, data interface{}) {
	t, err := template.ParseFiles(filename)
	if err != nil {
		errorHandler(err)
		return
	}
	err = t.Execute(responseWriter, data)
	if err != nil {
		errorHandler(err)
		return
	}
}

func setCookie(responseWriter http.ResponseWriter, name string, value string) {
	cookieToStore := http.Cookie{Name: name, Value: value}
	http.SetCookie(responseWriter, &cookieToStore)
}

func errorHandler(err error) {
	if err != nil {
		fmt.Println(err)
		panic(err.Error())
	}
}

func closeDB(db *sql.DB) {
	err := db.Close()
	errorHandler(err)
}

func openDB() *sql.DB {
	db, err := sql.Open("mysql", "root:admin@tcp(127.0.0.1:3306)/todoapi")
	errorHandler(err)
	return db
}
