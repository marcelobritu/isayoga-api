package handler

import (
	"encoding/json"
	"net/http"

	"github.com/marcelobritu/isayoga-api/internal/usecase/class"
	"github.com/marcelobritu/isayoga-api/pkg/logger"
	"go.uber.org/zap"
)

type ClassHandler struct {
	createClass *class.CreateClassUseCase
	listClasses *class.ListClassesUseCase
}

func NewClassHandler(
	createClass *class.CreateClassUseCase,
	listClasses *class.ListClassesUseCase,
) *ClassHandler {
	return &ClassHandler{
		createClass: createClass,
		listClasses: listClasses,
	}
}

func (h *ClassHandler) Create(w http.ResponseWriter, r *http.Request) {
	var input class.CreateClassInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		logger.Error("Erro ao decodificar requisição", zap.Error(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	class, err := h.createClass.Execute(r.Context(), input)
	if err != nil {
		logger.Error("Erro ao criar aula", zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(class)
}

func (h *ClassHandler) List(w http.ResponseWriter, r *http.Request) {
	classes, err := h.listClasses.Execute(r.Context())
	if err != nil {
		logger.Error("Erro ao listar aulas", zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(classes)
}

