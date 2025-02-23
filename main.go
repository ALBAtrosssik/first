package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
)

type requestBody struct {
	Message string `json:"message"`
} // структура которая парсит json в go

func TaskHandler(w http.ResponseWriter, r *http.Request) {
	var taskReq requestBody                         // переменная для хранения декодированных данных
	err := json.NewDecoder(r.Body).Decode(&taskReq) // создание нового декодера который читает тело запроса

	if err != nil || taskReq.Message == "" {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}
	
	newTask := Task{Task: taskReq.Message, IsDone: false}
	data := DB.Create(&newTask)

	if data.Error != nil {
		json.NewEncoder(w).Encode(map[string]interface{}{"message": data.Error})
	} else {
		json.NewEncoder(w).Encode(map[string]interface{}{"message": newTask.ID})
	}
}

func AllHandler(w http.ResponseWriter, r *http.Request) {
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

func main() {
	// Вызываем метод InitDB() из файла db.go
	InitDB()
	// Автоматическая миграция модели Task
	DB.AutoMigrate(&Task{})

	router := mux.NewRouter()
	router.HandleFunc("/api/alltask", AllHandler).Methods("GET")
	router.HandleFunc("/api/task", TaskHandler).Methods("POST")
	http.ListenAndServe(":8080", router)
}
