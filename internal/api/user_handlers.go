package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/barretot/gobid/internal/jsonutils"
	"github.com/barretot/gobid/internal/services"
	"github.com/barretot/gobid/internal/usecase/user"
	"github.com/barretot/gobid/internal/validator"
)

func (api *Api) handleSignupUser(w http.ResponseWriter, r *http.Request) {
	var req user.CreateUserReq

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Corpo da requisição inválido", http.StatusBadRequest)
		return
	}

	if err := validator.ValidateRequest(req); err != nil {
		http.Error(w, fmt.Sprintf("Erro de validação: %v", err), http.StatusBadRequest)
		return
	}

	id, err := api.UserService.CreateUser(
		r.Context(),
		req.UserName,
		req.Email,
		req.Password,
		req.Bio,
	)

	if err != nil {
		errors.Is(err, services.ErrDuplicatedEmailOrUsername)
		_ = jsonutils.EncodeJson(w, r, http.StatusUnprocessableEntity, map[string]any{
			"error": "email or username already exists",
		})
		return
	}

	_ = jsonutils.EncodeJson(w, r, http.StatusUnprocessableEntity, map[string]any{
		"user_id": id,
	})
}

func (api *Api) handleLoginUser(w http.ResponseWriter, r *http.Request) {
	var req user.LoginUserReq

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Corpo da requisição inválido", http.StatusBadRequest)
		return
	}

	if err := validator.ValidateRequest(req); err != nil {
		http.Error(w, fmt.Sprintf("validation error: %v", err), http.StatusBadRequest)
		return
	}

	id, err := api.UserService.AuthenticateUser(r.Context(), req.Email, req.Password)

	if err != nil {
		if errors.Is(err, services.ErrInvalidCredentials) {
			jsonutils.EncodeJson(w, r, http.StatusUnauthorized, map[string]any{
				"error": "invalid email or password",
			})
			return
		}
		jsonutils.EncodeJson(w, r, http.StatusInternalServerError, map[string]any{
			"error": "unexpected internal server error",
		})
		return
	}

	err = api.Sessions.RenewToken(r.Context())
	if err != nil {
		jsonutils.EncodeJson(w, r, http.StatusInternalServerError, map[string]any{
			"error": "unexpected internal server error",
		})
		return
	}

	api.Sessions.Put(r.Context(), "AuthenticatedUserId", id)

	jsonutils.EncodeJson(w, r, http.StatusOK, map[string]any{
		"message": "logged in successfully",
	})
}

func (api *Api) handleLogoutUser(w http.ResponseWriter, r *http.Request) {
	err := api.Sessions.RenewToken(r.Context())
	if err != nil {
		jsonutils.EncodeJson(w, r, http.StatusInternalServerError, map[string]any{
			"error": "unexpected internal server error",
		})
		return
	}

	api.Sessions.Remove(r.Context(), "AuthenticatedUserId")
	jsonutils.EncodeJson(w, r, http.StatusOK, map[string]any{
		"message": "logged out successfully",
	})
}
