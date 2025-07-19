package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

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
		if errors.Is(err, services.ErrDuplicatedEmailOrPassword) {
			http.Error(w, "Email ou nome de usuário já está em uso", http.StatusConflict)
			return
		}
		http.Error(w, "Erro interno ao criar usuário", http.StatusInternalServerError)
		return
	}

	// Sucesso
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, `{"id":"%s","message":"Usuário criado com sucesso"}`, id.String())
}

func (api *Api) handleLoginUser(w http.ResponseWriter, r *http.Request) {
	panic("TODO - NOT IMPLEMENTED")
}

func (api *Api) handleLogoutUser(w http.ResponseWriter, r *http.Request) {
	panic("TODO - NOT IMPLEMENTED")
}
