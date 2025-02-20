package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

var task string // объявление глобальной переменной

type requestBody struct {
	Message string `json:"message"`
} // структура которая парсит json в go

func TaskHandler(w http.ResponseWriter, r *http.Request) {
	var taskReq requestBody // переменная для хранения декодированных данных

	err := json.NewDecoder(r.Body).Decode(&taskReq) // создание нового декодера который читает тело запроса
	// метод decode парсит json и заполняет taskReq
	if err != nil || taskReq.Message == "" {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	} //проверка, если есть ошибка или пустое сообщение

	task = taskReq.Message
	w.WriteHeader(http.StatusOK)                                                                  // отправляет серверу статус ок
	json.NewEncoder(w).Encode(map[string]string{"status": "success", "message": "Task received"}) // отправляем json ответ
}

func HelloHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json") // пишем для того чтобы клиент понимал что данные в формате json

	greeting := "Hello"
	if task != "" {
		greeting = fmt.Sprintf("Hello, %s", task)
	}
	// если task не пустой отправляется приветвие с task
	response := map[string]string{"greeting": greeting} // создаем мапу с ключем greeting
	json.NewEncoder(w).Encode(response)                 // кодируем сообщение в json
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/api/hello", HelloHandler).Methods("GET")
	router.HandleFunc("/api/task", TaskHandler).Methods("POST")

	http.ListenAndServe(":8080", router)
}
