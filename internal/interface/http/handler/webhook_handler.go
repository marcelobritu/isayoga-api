package handler

import (
	"encoding/json"
	"net/http"

	"github.com/marcelobritu/isayoga-api/internal/usecase/payment"
	"github.com/marcelobritu/isayoga-api/pkg/logger"
	"go.uber.org/zap"
)

type WebhookHandler struct {
	processWebhook *payment.ProcessWebhookUseCase
}

func NewWebhookHandler(processWebhook *payment.ProcessWebhookUseCase) *WebhookHandler {
	return &WebhookHandler{
		processWebhook: processWebhook,
	}
}

func (h *WebhookHandler) MercadoPago(w http.ResponseWriter, r *http.Request) {
	var input payment.WebhookInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		logger.Error("Erro ao decodificar webhook", zap.Error(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	logger.Info("Webhook recebido",
		zap.String("type", input.Type),
		zap.String("action", input.Action),
	)

	if err := h.processWebhook.Execute(r.Context(), input); err != nil {
		logger.Error("Erro ao processar webhook", zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

