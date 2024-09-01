package entity

import (
	"time"
)

type CreateProductRequest struct {
	ShopId string `validate:"uuid" db:"shop_id"`

	Name        string `json:"name" validate:"required" db:"name"`
	Description string `json:"description" validate:"required,max=255" db:"description"`
	Price       int    `json:"price" validate:"required,min=0" db:"price"`
	Stock       int    `json:"stock" validate:"required,min=0" db:"stock"`
	Category    string `json:"category" validate:"required" db:"category"`
}

type CreateProductResponse struct {
	Id string `json:"id" db:"id"`
}

type GetAllProductRequest struct {
	ShopId string `validate:"uuid" db:"shop_id"`
	Page   int    `json:"page" validate:"min=1"`
	Limit  int    `json:"limit" validate:"min=1,max=100"`
}

type GetAllProductResponse struct {
	Products []GetProductResponse `json:"products"`
	Total    int                  `json:"total"`
	Page     int                  `json:"page"`
	Limit    int                  `json:"limit"`
}

type GetProductRequest struct {
	Id string `validate:"uuid" db:"id"`
}

type GetProductResponse struct {
	Name        string `json:"name" db:"name"`
	Description string `json:"description" db:"description"`
	Price       int    `json:"price" db:"price"`
	Stock       int    `json:"stock" db:"stock"`
	Category    string `json:"category" db:"category"`
}

type UpdateProductRequest struct {
	Id string `validate:"uuid" db:"id"`

	Name        string `json:"name" validate:"required" db:"name"`
	Description string `json:"description" validate:"required,max=255" db:"description"`
	Price       int    `json:"price" validate:"required,min=0" db:"price"`
	Stock       int    `json:"stock" validate:"required,min=0" db:"stock"`
	Category    string `json:"category" validate:"required" db:"category"`
}

type UpdateProductResponse struct {
	Id string `json:"id" db:"id"`

	Name        string `json:"name" db:"name"`
	Description string `json:"description" db:"description"`
	Price       int    `json:"price" db:"price"`
	Stock       int    `json:"stock" db:"stock"`
	Category    string `json:"category" db:"category"`
}

type DeleteProductRequest struct {
	Id string `validate:"uuid" db:"id"`
}

type DeleteProductResponse struct {
	Id string `json:"id" db:"id"`
}

type Product struct {
	Id          string     `json:"id" db:"id"`
	ShopId      string     `json:"shop_id" db:"shop_id"`
	Name        string     `json:"name" db:"name"`
	Description string     `json:"description" db:"description"`
	Price       int        `json:"price" db:"price"`
	Stock       int        `json:"stock" db:"stock"`
	CategoryId  *string    `json:"category_id" db:"category_id"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty" db:"deleted_at"`
}

type ProductsRequest struct {
	Page     int `query:"page" validate:"required"`
	Paginate int `query:"paginate" validate:"required"`
}

func (r *ProductsRequest) SetDefault() {
	if r.Page < 1 {
		r.Page = 1
	}

	if r.Paginate < 1 {
		r.Paginate = 10
	}
}
