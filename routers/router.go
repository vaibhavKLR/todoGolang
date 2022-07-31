package router

import (
	"github.com/gorilla/mux"
	controller "github.com/vaibhavKLR/todoApp/controllers"
)

func Router() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/api/todos", controller.GetAllTodos).Methods("GET")
	router.HandleFunc("/api/todos/{id}", controller.GetTodoById).Methods("GET")
	router.HandleFunc("/api/todo", controller.CreateTodo).Methods("POST")
	router.HandleFunc("/api/todo/{id}", controller.MarkAsDone).Methods("PUT")
	router.HandleFunc("/api/todo/undone/{id}", controller.MarkAsUndone).Methods("PUT")
	router.HandleFunc("/api/todo/{id}", controller.DeleteATodo).Methods("DELETE")
	router.HandleFunc("/api/todo/edit/{id}", controller.EditOneTodo).Methods("PUT")

	return router
}
