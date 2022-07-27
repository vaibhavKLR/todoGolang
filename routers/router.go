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
	router.HandleFunc("/api/todo/{id}", controller.DeleteATodo).Methods("DELETE")

	return router
}
