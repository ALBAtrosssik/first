package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
)

var task Task

type requestBody struct {
	Message string `json:"message"`
}

func TaskHandler(w http.ResponseWriter, r *http.Request) {
	var taskReq requestBody
	err := json.NewDecoder(r.Body).Decode(&taskReq)

	if err != nil || taskReq.Message == "" {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	newTask := Task{Task: taskReq.Message, IsDone: false}
	data := DB.Create(&newTask)

	if data.Error != nil {
		http.Error(w, data.Error.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(newTask)
}

func TasksHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var tasks []Task
	result := DB.Find(&tasks)

	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(tasks); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func UpdateHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	result := DB.First(&task, id)

	if result.Error != nil {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}

	if task.IsDone == false {
		DB.Model(&task).Update("is_done", true)
	} else {
		DB.Model(&task).Update("is_done", false)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"id": id, "is_done": task.IsDone})
}

func DeleteHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	result := DB.First(&task, id)

	if result.Error != nil {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}

	DB.Delete(&task)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"id": id})
}

func main() {
	// Вызываем метод InitDB() из файла db.go
	InitDB()
	// Автоматическая миграция модели Task
	DB.AutoMigrate(&Task{})

	router := mux.NewRouter()
	router.HandleFunc("/api/tasks", TasksHandler).Methods("GET")
	router.HandleFunc("/api/task", TaskHandler).Methods("POST")
	router.HandleFunc("/api/task/{id}", UpdateHandler).Methods("PATCH")
	router.HandleFunc("/api/task/{id}", DeleteHandler).Methods("DELETE")
	http.ListenAndServe(":8080", router)
}
