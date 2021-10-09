package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"./models"
)

var dsn = "root:root@tcp(go-todolist-api_devcontainer_mysql_1:3306)/todolist?charset=utf8mb4&parseTime=True&loc=Local"
var db, _ = gorm.Open(mysql.Open(dsn), &gorm.Config{})

func init() {
	log.SetFormatter(&log.TextFormatter{})
	log.SetReportCaller(true)
}

func main() {
	db.Debug().AutoMigrate(&models.Todo{})

	log.Info("Starting Todolist API server")
	router := mux.NewRouter()
	router.HandleFunc("/heartbeat", Heartbeat).Methods("GET")
	router.HandleFunc("/todo-complete", GetCompleted).Methods("GET")
	router.HandleFunc("/todo-incomplete", GetIncomplete).Methods("GET")
	router.HandleFunc("/todo", Get).Methods("GET")
	router.HandleFunc("/todo", Create).Methods("POST")
	router.HandleFunc("/todo/{id}", Update).Methods("POST", "PUT")
	router.HandleFunc("/todo/{id}", Delete).Methods("DELETE")

	handler := cors.New(cors.Options{
		AllowedMethods: []string{"GET", "POST", "DELETE", "PATCH", "OPTIONS", "PUT"},
	}).Handler(router)

	http.ListenAndServe(":3010", handler)
}

func Heartbeat(w http.ResponseWriter, r *http.Request) {
	log.Info("Heartbeat called")
	w.Header().Set("Content-Type", "application/json")
	response := models.Response{
		Success: true,
		Payload: "Service is up.",
		Error:   "",
	}
	json.NewEncoder(w).Encode(response)
}

func Create(w http.ResponseWriter, r *http.Request) {
	request := models.CreateRequest{}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	log.WithFields(log.Fields{"description": request.Description}).Info("Add new Todo. Saving to database.")
	todo := models.Todo{}
	todo.Description = request.Description
	db.Create(&todo)
	db.Last(&todo)
	w.Header().Set("Content-Type", "application/json")
	response := models.Response{
		Success: true,
		Payload: strconv.Itoa(todo.Id),
		Error:   "",
	}
	json.NewEncoder(w).Encode(response)
}

func Update(w http.ResponseWriter, r *http.Request) {
	// Get URL parameter from mux
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	request := models.UpdateRequest{}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Test if the TodoItem exist in DB
	todo := &models.Todo{}
	response := models.Response{}
	result := db.First(&todo, id)
	if result.Error != nil {
		w.Header().Set("Content-Type", "application/json")
		response.Error = result.Error.Error()
		response.Success = false
		json.NewEncoder(w).Encode(response)
	} else {
		log.WithFields(log.Fields{"Id": id, "Completed": request.IsCompleted}).Info("Updating Todo")
		db.First(&todo, id)
		todo.Completed = request.IsCompleted
		db.Save(&todo)
		w.Header().Set("Content-Type", "application/json")
		response.Success = true
		response.Payload = "Update successful"
		json.NewEncoder(w).Encode(response)
	}
}

func Delete(w http.ResponseWriter, r *http.Request) {
	// Get URL parameter from mux
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	// Test if the TodoItem exist in DB
	todo := models.Todo{}
	response := models.Response{}
	result := db.First(&todo, id)
	if result.Error != nil {
		w.Header().Set("Content-Type", "application/json")
		response.Error = result.Error.Error()
		response.Success = false
		json.NewEncoder(w).Encode(response)
	} else {
		log.WithFields(log.Fields{"Id": id}).Info("Deleting Todo")
		db.First(&todo, id)
		db.Delete(&todo)
		w.Header().Set("Content-Type", "application/json")
		response.Success = true
		response.Payload = "Delete successful"
		json.NewEncoder(w).Encode(response)
	}
}

func GetCompleted(w http.ResponseWriter, r *http.Request) {
	log.Info("Get completed TodoItems")
	completedTodoItems := GetTodoItems(true)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(completedTodoItems)
}

func GetIncomplete(w http.ResponseWriter, r *http.Request) {
	log.Info("Get Incomplete TodoItems")
	incompleteTodoItems := GetTodoItems(false)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(incompleteTodoItems)
}

func Get(w http.ResponseWriter, r *http.Request) {
	log.Info("Get All TodoItems")
	var todos []models.Todo
	db.Find(&todos)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&todos)
}

//Good example how gorm works
func GetTodoItems(completed bool) interface{} {
	var todos []models.Todo
	db.Where("completed = ?", completed).Find(&todos)
	return todos
}
