package handler

import (
	"context"
	"devbubble-api/internal/core"

	"devbubble-api/pkg/json"
	"devbubble-api/pkg/validator"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
)

type UserService interface {
	GetUserById(id int) string
	CreateUser(ctx context.Context, userDto core.CreateUserDto) (string, error)
	VerifyUser(ctx context.Context, dto core.VerifyUserDto) (string, error)
	GetCurrentUser(ctx context.Context, id string) (core.CurrentUserResponse, error)
	UpdateUser(ctx context.Context, userDto core.UpdateUserDto) error
}

type UserHandler struct {
	userService   UserService
	authService   AuthService
	jwtMiddleware func(next http.Handler) http.Handler
	validator     *validator.Validator
	log           *slog.Logger
}

func NewUserHandler(userService UserService, authService AuthService, validate *validator.Validator, log *slog.Logger, jwtMiddleware func(next http.Handler) http.Handler) *UserHandler {
	return &UserHandler{
		userService, authService, jwtMiddleware,
		validate, log,
	}
}

func (h *UserHandler) InitRouter() chi.Router {
	r := chi.NewRouter()

	r.Post("/signUp", h.createUser)
	r.Post("/verify", h.verifyUser)
	r.With(h.jwtMiddleware).Get("/current", h.getCurrentUser)
	r.With(h.jwtMiddleware).Patch("/", h.updateUser)
	return r
}

func (h *UserHandler) getCurrentUser(w http.ResponseWriter, r *http.Request) {
	id, ok := r.Context().Value("UserId").(string) // Type assertion to string
	if !ok {
		// If the assertion fails or the UserId is not found, return an error or handle it appropriately.
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	user, err := h.userService.GetCurrentUser(r.Context(), id)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	json.JsonResponse(w, http.StatusOK, user)
}

func (h *UserHandler) updateUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id, ok := ctx.Value("UserId").(string) // Type assertion to string
	if !ok {
		// If the assertion fails or the UserId is not found, return an error or handle it appropriately.
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	user := core.UpdateUserDto{Id: id}
	err := json.DecodeJSONRequest(r, &user)
	if err != nil {
		json.HttpError(w, http.StatusBadRequest, err.Error())
		return
	}

	err = h.validator.Validate(w, user)
	if err != nil {

		return
	}

	err = h.userService.UpdateUser(ctx, user)
	if err != nil {

		json.HttpError(w, http.StatusBadRequest, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *UserHandler) createUser(w http.ResponseWriter, r *http.Request) {
	user := core.CreateUserDto{}
	err := json.DecodeJSONRequest(r, &user)
	if err != nil {
		json.HttpError(w, http.StatusBadRequest, err.Error())
		return
	}

	err = h.validator.Validate(w, user)
	if err != nil {

		return
	}

	email, err := h.userService.CreateUser(r.Context(), user)
	if err != nil {
		json.HttpError(w, http.StatusBadRequest, err.Error())
		return
	}
	json.JsonResponse(w, http.StatusCreated, core.CreateUserResponse{Email: email})

}

func (h *UserHandler) verifyUser(w http.ResponseWriter, r *http.Request) {
	dto := core.VerifyUserDto{}
	err := json.DecodeJSONRequest(r, &dto)
	if err != nil {
		json.HttpError(w, http.StatusBadRequest, err.Error())
		return
	}
	err = h.validator.Validate(w, dto)
	if err != nil {

		return
	}
	id, err := h.userService.VerifyUser(r.Context(), dto)
	if err != nil {
		json.HttpError(w, http.StatusUnauthorized, "wrong code")
		return
	}

	token, err := h.authService.CreateJWT(id)
	if err != nil {
		h.log.Error("jwt Error")
		json.HttpError(w, http.StatusUnauthorized, err.Error())
		return
	}

	cookie := &http.Cookie{
		Name:     "Session",
		Value:    token,
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
		Path:     "/",
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	}

	http.SetCookie(w, cookie)

	w.WriteHeader(http.StatusOK)

}
