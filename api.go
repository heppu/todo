package todo

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

type Handler struct {
	todoList *List
	router   *httprouter.Router
}

func NewHandler(todoList *List) *Handler {
	handler := &Handler{
		todoList: todoList,
		router:   httprouter.New(),
	}

	handler.router.GET("/api/tasks/:id", handler.GetTask)
	handler.router.POST("/api/tasks", handler.PostTask)

	return handler
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.router.ServeHTTP(w, r)
}

func (h *Handler) GetTask(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	id, err := strconv.ParseUint(p.ByName("id"), 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Println(err)
		return
	}

	task, err := h.todoList.TaskByID(id)
	if err == ErrTaskNotFound {
		http.Error(w, err.Error(), http.StatusNotFound)
		log.Println(err)
		return
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println(err)
		return
	}

	if err := json.NewEncoder(w).Encode(task); err != nil {
		log.Println(err)
	}
}

func (h *Handler) PostTask(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	task := TaskData{}
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Println(err)
		return
	}

	id, err := h.todoList.Add(task)
	if err == ErrEmptyTaskName || err == ErrTooLongTaskName {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Println(err)
		return
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println(err)
		return
	}

	if err := json.NewEncoder(w).Encode(id); err != nil {
		log.Println(err)
	}
}
