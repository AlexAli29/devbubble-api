package handler

import (
	"context"
	"devbubble-api/internal/core"
	db "devbubble-api/internal/repository"
	"devbubble-api/pkg/json"
	"devbubble-api/pkg/validator"

	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type MessageService interface {
	CreateMessage(ctx context.Context, dto core.CreateMessageDto) (core.Message, error)
	RemoveMessage(ctx context.Context, dto core.RemoveMessageDto) error
	GetMessagesByChatId(ctx context.Context, dto core.GetMessagesByChatIdDto) ([]db.GetMessagesByChatIDRow, error)
}
type MessageHandler struct {
	messageService MessageService
	validator      *validator.Validator
	jwtMiddleware  func(next http.Handler) http.Handler
	log            *slog.Logger
}

func NewMessageHandler(messageService MessageService, validator *validator.Validator, jwtMiddleware func(next http.Handler) http.Handler, log *slog.Logger) *MessageHandler {

	return &MessageHandler{messageService, validator, jwtMiddleware, log}

}

func (h *MessageHandler) InitRouter() chi.Router {
	r := chi.NewRouter()
	r.With(h.jwtMiddleware).Post("/", h.createMessage)
	r.With(h.jwtMiddleware).Delete("/{id}", h.removeMessage)

	return r
}
func (h *MessageHandler) createMessage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userId, ok := ctx.Value("UserId").(string)
	if !ok {

		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	dto := core.CreateMessageDto{UserId: userId}
	err := json.DecodeJSONRequest(r, &dto)
	if err != nil {
		json.HttpError(w, http.StatusBadRequest, err.Error())
		return
	}

	err = h.validator.Validate(w, dto)
	if err != nil {

		return
	}

	message, err := h.messageService.CreateMessage(ctx, dto)
	if err != nil {
		json.HttpError(w, http.StatusBadRequest, err.Error())
		return
	}
	json.JsonResponse(w, http.StatusCreated, message)

}

func (h *MessageHandler) removeMessage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userId, ok := ctx.Value("UserId").(string)
	if !ok {

		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	messageId := chi.URLParam(r, "id")
	dto := core.RemoveMessageDto{UserId: userId, Id: messageId}

	err := h.validator.Validate(w, dto)
	if err != nil {

		return
	}

	err = h.messageService.RemoveMessage(ctx, dto)
	if err != nil {
		json.HttpError(w, http.StatusBadRequest, err.Error())
		return
	}
	w.WriteHeader(http.StatusOK)

}
