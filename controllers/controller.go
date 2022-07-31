package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	model "github.com/vaibhavKLR/todoApp/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	mongo_bson "gopkg.in/mgo.v2/bson"
)

const connectionString = "mongodb+srv://vaibhavs828:vaibhav@cluster0.ngtma.mongodb.net/?retryWrites=true&w=majority"
const dbName = "todoApp"
const collName = "todos"

//Reference
var collection *mongo.Collection

//Connection
func init() {

	// Set client options
	clientOption := options.Client().ApplyURI(connectionString)

	// connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOption)

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("MongoDB connection success")

	collection = client.Database(dbName).Collection(collName)

	fmt.Println("Collection reference ready")

}

func insertOneTodo(todo model.Todo) bson.M{
	inserted, err := collection.InsertOne(context.Background(), todo)

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("The Id is ", inserted.InsertedID)

	var insertedTodo bson.M
	cursor := collection.FindOne(context.Background(), bson.M{"_id": inserted.InsertedID}).Decode(&insertedTodo)
	if cursor!=nil {
		log.Fatal(cursor)
	}
	return insertedTodo
}

//updating a todo if completed
func updateOneTodo(todoID string) bool {

	if !mongo_bson.IsObjectIdHex(todoID) {
		return false
	}
	id, _ := primitive.ObjectIDFromHex(todoID)
	filter := bson.M{"_id": id}
	update := bson.M{"$set": bson.M{"completed": true}}
	result, err := collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Modified count is ", result.ModifiedCount)
	return true
}

func updateAsUndone(todoID string) bool {
	if !mongo_bson.IsObjectIdHex(todoID) {
		return false
	}
	id, _ := primitive.ObjectIDFromHex(todoID)
	filter := bson.M{"_id": id}
	update := bson.M{"$set": bson.M{"completed": false}}
	result, err := collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Modified count is ", result.ModifiedCount)
	return true
}

//Edit the task
func editOneTodo(todoID string, newTask model.Todo) bool{
	id, _ := primitive.ObjectIDFromHex(todoID)

	//If ID not found
	if !mongo_bson.IsObjectIdHex(todoID) {
		return false
	}

	filter := bson.M{"_id": id}
	//The line below is working but not working for body
	//update := bson.M{"$set": bson.M{"task": "Will be updated"}}
	update := bson.M{"$set": bson.M{"task": newTask.Task}}

	result, err := collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Modified count is ", result.ModifiedCount)
	//fmt.Println("Task is ")
	return true
}

//delete one todo
func deleteOneTodo(todoID string) {
	id, _ := primitive.ObjectIDFromHex(todoID)
	filter := bson.M{"_id": id}
	deleteCount, err := collection.DeleteOne(context.Background(), filter)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("The delete count is ", deleteCount)
}

func getTodoById(todoID string) (bool, bson.M) {

	if !mongo_bson.IsObjectIdHex(todoID) {
		return false, bson.M{}
	}

	id, _ := primitive.ObjectIDFromHex(todoID)
	var todo bson.M
	cursor := collection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&todo)

	//This terminates the connection
	if cursor != nil {
		log.Fatal(cursor)
	}

	return true, todo

	//fmt.Println("Result is ", cursor)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	var todos []primitive.M
	// 	var todo bson.M
	// 	err1 := result.Decode(&todo)
	// 	if err1 != nil {
	// 		log.Fatal(err)
	// 	}
	// 	todos = append(todos, todo)
	// 	return todos
}

//get all todos
func getAllTodos() []primitive.M {
	cursor, err := collection.Find(context.Background(), bson.D{{}})
	//returns a cursor, {{}} means selecting all

	if err != nil {
		log.Fatal(err)
	}

	var todos []primitive.M
	for cursor.Next(context.Background()) {
		var todo bson.M
		err := cursor.Decode(&todo)
		if err != nil {
			log.Fatal(err)
		}
		todos = append(todos, todo)
	}
	defer cursor.Close(context.Background())
	return todos
}

//Controllers - Functions to be called

func GetTodoById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-www-form-urlencode")
	params := mux.Vars(r)
	//fmt.Println("Inside GetTOdoById")
	err, todoById := getTodoById(params["id"])
	if !err {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Wrong ID"))
		return
	}
	json.NewEncoder(w).Encode(todoById)
}

func GetAllTodos(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-www-form-urlencode")
	alltodos := getAllTodos()
	json.NewEncoder(w).Encode(alltodos)
}

func CreateTodo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-www-form-urlencode")
	w.Header().Set("Allow-Control-Allow-Methods", "POST")

	var todo model.Todo
	_ = json.NewDecoder(r.Body).Decode(&todo)
	newTodo := insertOneTodo(todo)
	json.NewEncoder(w).Encode(newTodo)
}

func MarkAsDone(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-www-form-urlencode")
	w.Header().Set("Allow-Control-Allow-Methods", "PUT")

	params := mux.Vars(r)
	err := updateOneTodo(params["id"])
	//IF ID not found
	if !err {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Wrong ID"))
		return
	}
	json.NewEncoder(w).Encode(params["id"])
}

func MarkAsUndone(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-www-form-urlencode")
	w.Header().Set("Allow-Control-Allow-Methods", "PUT")
	params := mux.Vars(r)
	err := updateAsUndone(params["id"])
	//IF ID not found
	if !err {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Wrong ID"))
		return
	}
	json.NewEncoder(w).Encode(params["id"])

}

func EditOneTodo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-www-form-urlencode")
	w.Header().Set("Allow-Control-Allow-Methods", "PUT")

	params := mux.Vars(r)
	var updatedTaskTodo model.Todo
	_ = json.NewDecoder(r.Body).Decode(&updatedTaskTodo)
	err := editOneTodo(params["id"], updatedTaskTodo)
	if !err {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Wrong ID"))
		return
	}
	json.NewEncoder(w).Encode(updatedTaskTodo)
}

func DeleteATodo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-www-form-urlencode")
	w.Header().Set("Allow-Control-Allow-Methods", "DELETE")

	params := mux.Vars(r)
	deleteOneTodo(params["id"])
	json.NewEncoder(w).Encode(params["id"])
}
