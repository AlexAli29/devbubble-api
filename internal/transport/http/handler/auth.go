package handler

import (
	"context"
	"devbubble-api/internal/core"
	"devbubble-api/pkg/json"
	"devbubble-api/pkg/validator"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
)

type AuthService interface {
	CreateJWT(userId string) (string, error)
	SignIn(ctx context.Context, dto core.SignInDto) (string, error)
	ParseToken(tokenString string) (*core.AuthTokenClaims, error)
	GenerateAuthCode(ctx context.Context, email string) error
}
type AuthHandler struct {
	authService AuthService
	validate    *validator.Validator
}

func NewAuthHandler(authService AuthService, validate *validator.Validator) *AuthHandler {

	return &AuthHandler{authService, validate}

}

func (h *AuthHandler) InitRouter() chi.Router {
	r := chi.NewRouter()

	r.Post("/signIn", h.SignIn)
	r.Post("/generate", h.GenerateAuthCode)
	return r
}

func (h *AuthHandler) GenerateAuthCode(w http.ResponseWriter, r *http.Request) {
	dto := core.GenerateAuthCodeDto{}
	err := json.DecodeJSONRequest(r, &dto)
	if err != nil {
		json.HttpError(w, http.StatusBadRequest, err.Error())
		return
	}
	err = h.validate.Validate(w, dto)
	if err != nil {
		json.HttpError(w, http.StatusBadRequest, err.Error())
		return
	}
	err = h.authService.GenerateAuthCode(r.Context(), dto.Email)
	if err != nil {
		json.HttpError(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *AuthHandler) SignIn(w http.ResponseWriter, r *http.Request) {
	dto := core.SignInDto{}
	err := json.DecodeJSONRequest(r, &dto)
	if err != nil {
		json.HttpError(w, http.StatusBadRequest, err.Error())
		return
	}
	err = h.validate.Validate(w, dto)
	if err != nil {
		json.HttpError(w, http.StatusBadRequest, err.Error())
		return
	}

	userId, err := h.authService.SignIn(r.Context(), dto)
	if err != nil {
		json.HttpError(w, http.StatusUnauthorized, "wrong code or email")
		return
	}
	token, err := h.authService.CreateJWT(userId)
	if err != nil {
		json.HttpError(w, http.StatusUnauthorized, "wrong code or email")
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
