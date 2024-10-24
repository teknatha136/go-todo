package handlers

import (
	"encoding/json"
	"net/http"

	"go-todo-app/models"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

type TodoHandler struct {
	DB *gorm.DB
}

func RegisterTodoHandlers(r *mux.Router, db *gorm.DB) {
	handler := &TodoHandler{DB: db}

	r.HandleFunc("/todos", handler.GetTodos).Methods("GET")
	r.HandleFunc("/todos", handler.CreateTodo).Methods("POST")
	r.HandleFunc("/todos/{id}", handler.GetTodo).Methods("GET")
	r.HandleFunc("/todos/{id}", handler.UpdateTodo).Methods("PUT")
	r.HandleFunc("/todos/{id}", handler.DeleteTodo).Methods("DELETE")
}

func (h *TodoHandler) GetTodos(w http.ResponseWriter, r *http.Request) {
	var todos []models.Todo
	h.DB.Find(&todos)
	json.NewEncoder(w).Encode(todos)
}

func (h *TodoHandler) CreateTodo(w http.ResponseWriter, r *http.Request) {
	var todo models.Todo
	json.NewDecoder(r.Body).Decode(&todo)
	h.DB.Create(&todo)
	json.NewEncoder(w).Encode(todo)
}

func (h *TodoHandler) GetTodo(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var todo models.Todo
	if err := h.DB.First(&todo, params["id"]).Error; err != nil {
		http.NotFound(w, r)
		return
	}
	json.NewEncoder(w).Encode(todo)
}

// UpdateTodo updates an existing to-do item by its ID.
func UpdateTodo(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]

		var todo models.Todo
		if err := db.First(&todo, id).Error; err != nil {
			http.Error(w, "To-do not found", http.StatusNotFound)
			return
		}

		var updatedTodo models.Todo
		if err := json.NewDecoder(r.Body).Decode(&updatedTodo); err != nil {
			http.Error(w, "Invalid input", http.StatusBadRequest)
			return
		}

		// Update the fields if they are provided
		todo.Title = updatedTodo.Title
		todo.Completed = updatedTodo.Completed

		db.Save(&todo)
		json.NewEncoder(w).Encode(todo)
	}
}

func (h *TodoHandler) UpdateTodo(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var todo models.Todo
	if err := h.DB.First(&todo, params["id"]).Error; err != nil {
		http.NotFound(w, r)
		return
	}

	json.NewDecoder(r.Body).Decode(&todo)
	h.DB.Save(&todo)
	json.NewEncoder(w).Encode(todo)
}

func (h *TodoHandler) DeleteTodo(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var todo models.Todo
	if err := h.DB.Delete(&todo, params["id"]).Error; err != nil {
		http.NotFound(w, r)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
