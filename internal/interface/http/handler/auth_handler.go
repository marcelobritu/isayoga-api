package handler

import (
	"encoding/json"
	"net/http"

	"github.com/marcelobritu/isayoga-api/internal/usecase/auth"
	"github.com/marcelobritu/isayoga-api/pkg/logger"
	"go.uber.org/zap"
)

type AuthHandler struct {
	loginUseCase    *auth.LoginUseCase
	registerUseCase *auth.RegisterUseCase
}

func NewAuthHandler(loginUseCase *auth.LoginUseCase, registerUseCase *auth.RegisterUseCase) *AuthHandler {
	return &AuthHandler{
		loginUseCase:    loginUseCase,
		registerUseCase: registerUseCase,
	}
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var input auth.LoginInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		logger.Error("Erro ao decodificar request", zap.Error(err))
		http.Error(w, "Dados inválidos", http.StatusBadRequest)
		return
	}

	output, err := h.loginUseCase.Execute(r.Context(), input)
	if err != nil {
		logger.Error("Erro no login", zap.Error(err))
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(output)
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var input auth.RegisterInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		logger.Error("Erro ao decodificar request", zap.Error(err))
		http.Error(w, "Dados inválidos", http.StatusBadRequest)
		return
	}

	output, err := h.registerUseCase.Execute(r.Context(), input)
	if err != nil {
		logger.Error("Erro no registro", zap.Error(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(output)
}

