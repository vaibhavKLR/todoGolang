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

func insertOneTodo(todo model.Todo) {
	inserted, err := collection.InsertOne(context.Background(), todo)

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("The Id is ", inserted.InsertedID)
}

//updating a todo if completed
func updateOneTodo(todoID string) {
	id, _ := primitive.ObjectIDFromHex(todoID)
	filter := bson.M{"_id": id}
	update := bson.M{"$set": bson.M{"completed": true}}

	result, err := collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Modified count is ", result.ModifiedCount)
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

func getTodoById(todoID string) bson.M {
	id, _ := primitive.ObjectIDFromHex(todoID)
	var todo bson.M
	cursor := collection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&todo)
	if cursor != nil {
		log.Fatal(cursor)
	}
	return todo

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
	todoById := getTodoById(params["id"])
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
	insertOneTodo(todo)
	json.NewEncoder(w).Encode(todo)
}

func MarkAsDone(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-www-form-urlencode")
	w.Header().Set("Allow-Control-Allow-Methods", "PUT")

	params := mux.Vars(r)
	updateOneTodo(params["id"])
	json.NewEncoder(w).Encode(params["id "])
}

func DeleteATodo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-www-form-urlencode")
	w.Header().Set("Allow-Control-Allow-Methods", "DELETE")

	params := mux.Vars(r)
	deleteOneTodo(params["id"])
	json.NewEncoder(w).Encode(params["id"])
}
