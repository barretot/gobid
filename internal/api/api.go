package api

import (
	"github.com/barretot/gobid/internal/services"
	"github.com/go-chi/chi/v5"
)

type Api struct {
	Router      *chi.Mux
	UserService services.UserService
} // Cada método para cada rota
