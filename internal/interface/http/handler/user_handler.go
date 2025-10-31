package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/marcelobritu/isayoga-api/internal/usecase/user"
	"github.com/marcelobritu/isayoga-api/pkg/logger"
	"go.uber.org/zap"
)

type UserHandler struct {
	createUserUC *user.CreateUserUseCase
	getUserUC    *user.GetUserUseCase
	listUsersUC  *user.ListUsersUseCase
	updateUserUC *user.UpdateUserUseCase
	deleteUserUC *user.DeleteUserUseCase
}

func NewUserHandler(
	createUserUC *user.CreateUserUseCase,
	getUserUC *user.GetUserUseCase,
	listUsersUC *user.ListUsersUseCase,
	updateUserUC *user.UpdateUserUseCase,
	deleteUserUC *user.DeleteUserUseCase,
) *UserHandler {
	return &UserHandler{
		createUserUC: createUserUC,
		getUserUC:    getUserUC,
		listUsersUC:  listUsersUC,
		updateUserUC: updateUserUC,
		deleteUserUC: deleteUserUC,
	}
}

func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	var input user.CreateUserInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		logger.Warn("Falha ao decodificar dados do usuário", zap.Error(err))
		http.Error(w, "Dados inválidos", http.StatusBadRequest)
		return
	}

	createdUser, err := h.createUserUC.Execute(r.Context(), input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdUser)
}

func (h *UserHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	foundUser, err := h.getUserUC.Execute(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(foundUser)
}

func (h *UserHandler) List(w http.ResponseWriter, r *http.Request) {
	users, err := h.listUsersUC.Execute(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func (h *UserHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var input user.UpdateUserInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		logger.Warn("Falha ao decodificar dados do usuário", zap.Error(err))
		http.Error(w, "Dados inválidos", http.StatusBadRequest)
		return
	}

	input.ID = id

	updatedUser, err := h.updateUserUC.Execute(r.Context(), input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedUser)
}

func (h *UserHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if err := h.deleteUserUC.Execute(r.Context(), id); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
