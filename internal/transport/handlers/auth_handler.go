package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/shrin00/pen/internal/transport/httpx"
	"github.com/shrin00/pen/internal/usecase"
)

type AuthHandler struct {
	auth *usecase.AuthService
}

func NewAuthHandler(auth *usecase.AuthService) *AuthHandler {
	return &AuthHandler{
		auth: auth,
	}
}

type authRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type authUserResponse struct {
	ID    string `json:"id"`
	Email string `json:"email"`
}

func (a *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req authRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Println("failed to parse the json body")
		httpx.WriteJSON(w, http.StatusBadRequest, map[string]string{
			"error": "invalid json body",
		})
		return
	}

	user, err := a.auth.Register(r.Context(), req.Email, req.Password)
	if err != nil {
		if errors.Is(err, usecase.ErrInvalidEmail) || errors.Is(err, usecase.ErrWeakPassword) {
			httpx.WriteJSON(w, http.StatusBadRequest, map[string]string{
				"error": err.Error(),
			})
			return
		}

		if errors.Is(err, usecase.ErrEmailAlreadyExists) {
			log.Println("user already exists")
			httpx.WriteJSON(w, http.StatusConflict, map[string]string{
				"error": "user already exists",
			})
			return
		}

		log.Println("failed to register the user")
		httpx.WriteJSON(w, http.StatusInternalServerError, map[string]string{
			"error": "failed to register user",
		})
		return
	}

	httpx.WriteJSON(w, http.StatusCreated, authUserResponse{
		ID:    user.ID,
		Email: user.Email,
	})
}

func (a *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req authRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Println("invalid json body")
		httpx.WriteJSON(w, http.StatusBadRequest, map[string]string{
			"error": "invalid json body",
		})
		return
	}

	user, user_token, err := a.auth.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		if errors.Is(err, usecase.ErrInvalidCredentials) {
			log.Println("invalid credentials")
			httpx.WriteJSON(w, http.StatusUnauthorized, map[string]string{
				"error": "invalid credentials",
			})
			return
		}

		log.Println("failed to login user")
		httpx.WriteJSON(w, http.StatusInternalServerError, map[string]string{
			"error": "failed to login user",
		})
		return
	}

	login_cookie := http.Cookie{
		Name:     "session_token",
		Value:    user_token,
		Expires:  time.Now().Add(a.auth.SessionTTL),
		HttpOnly: true,
	}
	http.SetCookie(w, &login_cookie)
	httpx.WriteJSON(w, http.StatusOK, authUserResponse{
		ID:    user.ID,
		Email: user.Email,
	})
}
