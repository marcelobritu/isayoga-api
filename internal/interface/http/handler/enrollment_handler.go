package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/marcelobritu/isayoga-api/internal/usecase/enrollment"
	"github.com/marcelobritu/isayoga-api/pkg/logger"
	"go.uber.org/zap"
)

type EnrollmentHandler struct {
	enrollStudent    *enrollment.EnrollStudentUseCase
	cancelEnrollment *enrollment.CancelEnrollmentUseCase
}

func NewEnrollmentHandler(
	enrollStudent *enrollment.EnrollStudentUseCase,
	cancelEnrollment *enrollment.CancelEnrollmentUseCase,
) *EnrollmentHandler {
	return &EnrollmentHandler{
		enrollStudent:    enrollStudent,
		cancelEnrollment: cancelEnrollment,
	}
}

func (h *EnrollmentHandler) Enroll(w http.ResponseWriter, r *http.Request) {
	var input enrollment.EnrollStudentInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		logger.Error("Erro ao decodificar requisição", zap.Error(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result, err := h.enrollStudent.Execute(r.Context(), input)
	if err != nil {
		logger.Error("Erro ao realizar inscrição", zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(result)
}

func (h *EnrollmentHandler) Cancel(w http.ResponseWriter, r *http.Request) {
	enrollmentID := chi.URLParam(r, "id")

	if err := h.cancelEnrollment.Execute(r.Context(), enrollmentID); err != nil {
		logger.Error("Erro ao cancelar inscrição", zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
