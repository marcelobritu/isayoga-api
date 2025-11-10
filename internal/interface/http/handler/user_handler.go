package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	pkgAuth "github.com/marcelobritu/isayoga-api/pkg/auth"
	"github.com/marcelobritu/isayoga-api/internal/infrastructure/http/middleware"
	"github.com/marcelobritu/isayoga-api/internal/usecase/user"
	"github.com/marcelobritu/isayoga-api/pkg/logger"
	"go.uber.org/zap"
)

type UserHandler struct {
	createUserUseCase     *user.CreateUserUseCase
	getUserUseCase        *user.GetUserUseCase
	listUsersUseCase      *user.ListUsersUseCase
	updateUserUseCase     *user.UpdateUserUseCase
	deleteUserUseCase     *user.DeleteUserUseCase
	changePasswordUseCase *user.ChangePasswordUseCase
}

func NewUserHandler(
	createUserUseCase *user.CreateUserUseCase,
	getUserUseCase *user.GetUserUseCase,
	listUsersUseCase *user.ListUsersUseCase,
	updateUserUseCase *user.UpdateUserUseCase,
	deleteUserUseCase *user.DeleteUserUseCase,
	changePasswordUseCase *user.ChangePasswordUseCase,
) *UserHandler {
	return &UserHandler{
		createUserUseCase:     createUserUseCase,
		getUserUseCase:        getUserUseCase,
		listUsersUseCase:      listUsersUseCase,
		updateUserUseCase:     updateUserUseCase,
		deleteUserUseCase:     deleteUserUseCase,
		changePasswordUseCase: changePasswordUseCase,
	}
}

func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	var input user.CreateUserInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		logger.Error("Erro ao decodificar request", zap.Error(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result, err := h.createUserUseCase.Execute(r.Context(), input)
	if err != nil {
		logger.Error("Erro ao criar usuário", zap.Error(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(result)
}

func (h *UserHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	result, err := h.getUserUseCase.Execute(r.Context(), id)
	if err != nil {
		logger.Error("Erro ao buscar usuário", zap.Error(err))
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func (h *UserHandler) List(w http.ResponseWriter, r *http.Request) {
	result, err := h.listUsersUseCase.Execute(r.Context())
	if err != nil {
		logger.Error("Erro ao listar usuários", zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func (h *UserHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var input user.UpdateUserInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		logger.Error("Erro ao decodificar request", zap.Error(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	input.ID = id

	result, err := h.updateUserUseCase.Execute(r.Context(), input)
	if err != nil {
		logger.Error("Erro ao atualizar usuário", zap.Error(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func (h *UserHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if err := h.deleteUserUseCase.Execute(r.Context(), id); err != nil {
		logger.Error("Erro ao deletar usuário", zap.Error(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *UserHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	var input user.ChangePasswordInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		logger.Error("Erro ao decodificar request", zap.Error(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	claims, ok := r.Context().Value(middleware.UserClaimsKey).(*pkgAuth.Claims)
	if !ok {
		http.Error(w, "Usuário não autenticado", http.StatusUnauthorized)
		return
	}
	input.UserID = claims.UserID

	if err := h.changePasswordUseCase.Execute(r.Context(), input); err != nil {
		logger.Error("Erro ao alterar senha", zap.Error(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Senha alterada com sucesso",
	})
}
