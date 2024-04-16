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

type TagService interface {
	GetAvailableTags(ctx context.Context, userId string) ([]db.UserTag, error)
	AddTag(ctx context.Context, dto core.AddTagDto) error
	RemoveTag(ctx context.Context, dto core.RemoveTagDto) error
}

type TagHandler struct {
	tagService    TagService
	validate      *validator.Validator
	jwtMiddleware func(next http.Handler) http.Handler

	log *slog.Logger
}

func NewTagHandler(tagService TagService, validate *validator.Validator, jwtMiddleware func(next http.Handler) http.Handler, log *slog.Logger) *TagHandler {

	return &TagHandler{tagService, validate, jwtMiddleware, log}

}
func (h *TagHandler) InitRouter() chi.Router {
	r := chi.NewRouter()
	r.With(h.jwtMiddleware).Get("/", h.GetAvailableTags)
	r.With(h.jwtMiddleware).Post("/", h.AddTag)
	r.With(h.jwtMiddleware).Delete("/{tagId}", h.RemoveTag)
	return r
}

func (h *TagHandler) GetAvailableTags(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userId, ok := ctx.Value("UserId").(string)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	tags, err := h.tagService.GetAvailableTags(ctx, userId)
	if err != nil {
		json.HttpError(w, http.StatusBadRequest, err.Error())
		return
	}
	json.JsonResponse(w, http.StatusOK, tags)
}
func (h *TagHandler) AddTag(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userId, ok := ctx.Value("UserId").(string)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	dto := core.AddTagDto{}
	err := json.DecodeJSONRequest(r, &dto)

	if err != nil {
		json.HttpError(w, http.StatusBadRequest, err.Error())
		return
	}
	dto.UserId = userId
	err = h.validate.Validate(w, dto)
	if err != nil {
		return
	}

	err = h.tagService.AddTag(ctx, dto)
	if err != nil {
		json.HttpError(w, http.StatusBadRequest, err.Error())
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (h *TagHandler) RemoveTag(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userId, ok := ctx.Value("UserId").(string)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	tagId := chi.URLParam(r, "tagId")
	dto := core.RemoveTagDto{TagId: tagId, UserId: userId}

	err := h.validate.Validate(w, dto)
	if err != nil {
		return
	}

	err = h.tagService.RemoveTag(ctx, dto)
	if err != nil {
		json.HttpError(w, http.StatusBadRequest, err.Error())
		return
	}
	w.WriteHeader(http.StatusOK)
}
