package product

import (
	"time"
)

type CreateProductReq struct {
	ProductName string    `json:"product_name" validate:"required"`
	Description string    `json:"description" validate:"required"`
	Baseprice   float64   `json:"baseprice" validate:"required"`
	AuctionEnd  time.Time `json:"auction_end" validate:"required"`
}
