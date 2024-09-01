package repository

import (
	"codebase-app/internal/module/product/entity"
	"codebase-app/internal/module/product/ports"
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

var _ ports.ProductRepository = &productRepository{}

type productRepository struct {
	db *sqlx.DB
}

func NewProductRepository(db *sqlx.DB) *productRepository {
	return &productRepository{
		db: db,
	}
}

func (r *productRepository) CreateProduct(ctx context.Context, req *entity.CreateProductRequest) (*entity.CreateProductResponse, error) {
	var resp = new(entity.CreateProductResponse)
	query := `
		INSERT INTO products (shop_id, name, description, price, stock, category)
		VALUES (?, ?, ?, ?, ?, ?) RETURNING id
	`

	err := r.db.QueryRowContext(ctx, r.db.Rebind(query),
		req.ShopId,
		req.Name,
		req.Description,
		req.Price,
		req.Stock,
		req.Category).Scan(&resp.Id)
	if err != nil {
		log.Error().Err(err).Any("payload", req).Msg("repository::CreateProduct - Failed to create product")
		return nil, err
	}

	return resp, nil
}

func (r *productRepository) GetProduct(ctx context.Context, req *entity.GetProductRequest) (*entity.GetProductResponse, error) {
	var resp = new(entity.GetProductResponse)
	query := `
		SELECT name, description, price, stock, category
		FROM products
		WHERE id = ? AND deleted_at IS NULL
	`

	err := r.db.QueryRowxContext(ctx, r.db.Rebind(query), req.Id).StructScan(resp)
	if err != nil {
		log.Error().Err(err).Any("payload", req).Msg("repository::GetProduct - Failed to get product")
		return nil, err
	}

	return resp, nil
}

func (r *productRepository) UpdateProduct(ctx context.Context, req *entity.UpdateProductRequest) (*entity.UpdateProductResponse, error) {
	var resp = new(entity.UpdateProductResponse)
	query := `
		UPDATE products
		SET name = ?, description = ?, price = ?, stock = ?, category = ?, updated_at = NOW()
		WHERE id = ? AND deleted_at IS NULL
		RETURNING id, name, description, price, stock, category
	`

	err := r.db.QueryRowxContext(ctx, r.db.Rebind(query),
		req.Name,
		req.Description,
		req.Price,
		req.Stock,
		req.Category,
		req.Id).Scan(&resp.Id, &resp.Name, &resp.Description, &resp.Price, &resp.Stock, &resp.Category)
	if err != nil {
		log.Error().Err(err).Any("payload", req).Msg("repository::UpdateProduct - Failed to update product")
		return nil, err
	}

	return resp, nil
}

func (r *productRepository) DeleteProduct(ctx context.Context, req *entity.DeleteProductRequest) error {
	query := `
		UPDATE products
		SET deleted_at = NOW()
		WHERE id = ? AND deleted_at IS NULL
	`

	_, err := r.db.ExecContext(ctx, r.db.Rebind(query), req.Id)
	if err != nil {
		log.Error().Err(err).Any("payload", req).Msg("repository::DeleteProduct - Failed to delete product")
		return err
	}
	return nil
}

func (r *productRepository) GetAllProducts(ctx context.Context, req *entity.GetAllProductRequest) (*entity.GetAllProductResponse, error) {
	var resp = new(entity.GetAllProductResponse)
	query := `
		SELECT id, name, description, price, stock, category
		FROM products
		WHERE shop_id = ? AND deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`

	offset := (req.Page - 1) * req.Limit

	rows, err := r.db.QueryxContext(ctx, r.db.Rebind(query), req.ShopId, req.Limit, offset)
	if err != nil {
		log.Error().Err(err).Any("payload", req).Msg("repository::GetAllProducts - Failed to get products")
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var product entity.GetProductResponse
		if err := rows.StructScan(&product); err != nil {
			log.Error().Err(err).Msg("repository::GetAllProducts - Failed to scan product")
			return nil, err
		}
		resp.Products = append(resp.Products, product)
	}

	// Get total count
	countQuery := `
		SELECT COUNT(*) 
		FROM products 
		WHERE shop_id = ? AND deleted_at IS NULL
	`
	err = r.db.QueryRowContext(ctx, r.db.Rebind(countQuery), req.ShopId).Scan(&resp.Total)
	if err != nil {
		log.Error().Err(err).Msg("repository::GetAllProducts - Failed to get total count")
		return nil, err
	}

	resp.Page = req.Page
	resp.Limit = req.Limit

	return resp, nil
}
