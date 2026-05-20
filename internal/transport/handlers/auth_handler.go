package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

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
		httpx.WriteJson(w, http.StatusBadRequest, map[string]string{
			"error": "invalid json body",
		})
	}

	user, err := a.auth.Register(r.Context(), req.Email, req.Password)
	if err != nil {
		if errors.Is(err, usecase.ErrEmailAlreadyExists) {
			log.Println("user already exists")
			httpx.WriteJson(w, http.StatusConflict, map[string]string{
				"error": "user already exists",
			})
			return
		}

		log.Println("failed to register the user")
		httpx.WriteJson(w, http.StatusInternalServerError, map[string]string{
			"error": "failed to register user",
		})
		return
	}

	httpx.WriteJson(w, http.StatusOK, authUserResponse{
		ID:    user.ID,
		Email: user.Email,
	})
}

func (a *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req authRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Println("invalid json body")
		httpx.WriteJson(w, http.StatusBadRequest, map[string]string{
			"error": "invalid json body",
		})
	}

	user, err := a.auth.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		if errors.Is(err, usecase.ErrInvalidCredentials) {
			log.Println("invalid credentials")
			httpx.WriteJson(w, http.StatusUnauthorized, map[string]string{
				"error": "invalid credentials",
			})
			return
		}

		log.Println("failed to login user")
		httpx.WriteJson(w, http.StatusInternalServerError, map[string]string{
			"error": "failed to login user",
		})
		return
	}

	httpx.WriteJson(w, http.StatusOK, authUserResponse{
		ID:    user.ID,
		Email: user.Email,
	})
}
