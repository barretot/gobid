package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/barretot/gobid/internal/jsonutils"
	"github.com/barretot/gobid/internal/services"
	"github.com/barretot/gobid/internal/usecase/product"

	"github.com/barretot/gobid/internal/validator"
	"github.com/google/uuid"
)

func (api *Api) handleCreateProduct(w http.ResponseWriter, r *http.Request) {
	var req product.CreateProductReq

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if err := validator.ValidateRequest(req); err != nil {
		http.Error(w, fmt.Sprintf("validation error: %v", err), http.StatusBadRequest)
		return
	}

	userId, ok := api.Sessions.Get(r.Context(), "AuthenticatedUserId").(uuid.UUID)

	if !ok {
		jsonutils.EncodeJson(w, r, http.StatusInternalServerError, map[string]any{
			"error": "unexpected error, try again later",
		})

		return
	}

	id, err := api.ProductService.CreateProduct(
		r.Context(),
		userId,
		req.ProductName,
		req.Description,
		req.Baseprice,
		req.AuctionEnd,
	)

	if err != nil {
		errors.Is(err, services.ErrDuplicatedEmailOrUsername)
		_ = jsonutils.EncodeJson(w, r, http.StatusInternalServerError, map[string]any{
			"error": "failed to create product auction try again later",
		})
		return
	}

	_ = jsonutils.EncodeJson(w, r, http.StatusCreated, map[string]any{
		"message":    "product created with success",
		"product_id": id,
	})
}
