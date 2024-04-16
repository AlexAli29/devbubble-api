package handler

import (
	"context"
	"devbubble-api/internal/core"
	"devbubble-api/pkg/json"
	"devbubble-api/pkg/validator"

	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type PrivateChatService interface {
	CreatePrivateChat(ctx context.Context, dto core.CreatePrivateChatDto) (string, error)
	RemovePrivateChat(ctx context.Context, dto core.RemovePrivateChatDto) error
	GetPrivateChats(ctx context.Context, userId string) ([]core.GetPrivateChatsResponse, error)
}

type ChatHandler struct {
	privateChatService PrivateChatService

	validator      *validator.Validator
	jwtMiddleware  func(next http.Handler) http.Handler
	log            *slog.Logger
	messageService MessageService
}

func NewChatHandler(privateChatService PrivateChatService, validator *validator.Validator, jwtMiddleware func(next http.Handler) http.Handler, log *slog.Logger, messageService MessageService) *ChatHandler {

	return &ChatHandler{privateChatService, validator, jwtMiddleware, log, messageService}

}

func (h *ChatHandler) InitRouter() chi.Router {
	r := chi.NewRouter()
	r.With(h.jwtMiddleware).Post("/", h.createPrivateChat)
	r.With(h.jwtMiddleware).Get("/", h.getPrivateChats)
	r.With(h.jwtMiddleware).Delete("/{id}", h.removePrivateChat)
	r.With(h.jwtMiddleware).Get("/messages/{id}", h.getChatMessages)

	return r
}

func (h *ChatHandler) createPrivateChat(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userId, ok := ctx.Value("UserId").(string)
	if !ok {

		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	chat := core.CreatePrivateChatDto{UserId: userId}
	err := json.DecodeJSONRequest(r, &chat)
	if err != nil {
		json.HttpError(w, http.StatusBadRequest, err.Error())
		return
	}

	err = h.validator.Validate(w, chat)
	if err != nil {

		return
	}

	id, err := h.privateChatService.CreatePrivateChat(ctx, chat)
	if err != nil {
		json.HttpError(w, http.StatusBadRequest, err.Error())
		return
	}
	json.JsonResponse(w, http.StatusCreated, core.CreatePrivateChatResponse{Id: id})

}

func (h *ChatHandler) removePrivateChat(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userId, ok := ctx.Value("UserId").(string)
	if !ok {

		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	chatId := chi.URLParam(r, "id")
	dto := core.RemovePrivateChatDto{UserId: userId, ChatId: chatId}

	err := h.validator.Validate(w, dto)
	if err != nil {

		return
	}

	err = h.privateChatService.RemovePrivateChat(ctx, dto)
	if err != nil {
		json.HttpError(w, http.StatusBadRequest, err.Error())
		return
	}
	w.WriteHeader(http.StatusOK)

}

func (h *ChatHandler) getPrivateChats(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userId, ok := ctx.Value("UserId").(string)
	if !ok {

		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	chats, err := h.privateChatService.GetPrivateChats(ctx, userId)
	if err != nil {
		json.HttpError(w, http.StatusBadRequest, err.Error())
		return
	}

	json.JsonResponse(w, http.StatusOK, chats)

}

func (h *ChatHandler) getChatMessages(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userId, ok := ctx.Value("UserId").(string)
	if !ok {

		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	chatId := chi.URLParam(r, "id")
	dto := core.GetMessagesByChatIdDto{
		UserId: userId,
		ChatId: chatId,
	}
	messages, err := h.messageService.GetMessagesByChatId(ctx, dto)
	if err != nil {
		json.HttpError(w, http.StatusBadRequest, err.Error())
		return
	}

	json.JsonResponse(w, http.StatusOK, messages)

}
