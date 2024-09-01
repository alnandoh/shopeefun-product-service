package handler

import (
	"codebase-app/internal/adapter"
	"codebase-app/internal/middleware"
	"codebase-app/internal/module/product/entity"
	"codebase-app/internal/module/product/ports"
	"codebase-app/internal/module/product/repository"
	"codebase-app/internal/module/product/service"
	"codebase-app/pkg/response"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

type productHandler struct {
	service ports.ProductService
}

func NewProductHandler() *productHandler {
	var (
		handler = new(productHandler)
		repo    = repository.NewProductRepository(adapter.Adapters.ShopeefunPostgres)
		service = service.NewProductService(repo)
	)
	handler.service = service

	return handler
}

func (h *productHandler) Register(router fiber.Router) {
	router.Get("/products", middleware.UserIdHeader, h.GetAllProducts)
	router.Post("/products", middleware.UserIdHeader, h.CreateProduct)
	router.Get("/products/:id", h.GetProduct)
	router.Delete("/products/:id", middleware.UserIdHeader, h.DeleteProduct)
	router.Patch("/products/:id", middleware.UserIdHeader, h.UpdateProduct)
}

func (h *productHandler) CreateProduct(c *fiber.Ctx) error {
	var req entity.CreateProductRequest
	if err := c.BodyParser(&req); err != nil {
		log.Warn().Err(err).Msg("handler::CreateProduct - Failed to parse request body")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error(err))
	}

	resp, err := h.service.CreateProduct(c.Context(), &req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.Error(err))
	}

	return c.Status(fiber.StatusCreated).JSON(response.Success(resp, ""))
}

func (h *productHandler) GetProduct(c *fiber.Ctx) error {
	var req entity.GetProductRequest
	req.Id = c.Params("id")

	resp, err := h.service.GetProduct(c.Context(), &req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.Error(err))
	}

	return c.Status(fiber.StatusOK).JSON(response.Success(resp, ""))
}

func (h *productHandler) UpdateProduct(c *fiber.Ctx) error {
	var req entity.UpdateProductRequest
	req.Id = c.Params("id")

	if err := c.BodyParser(&req); err != nil {
		log.Warn().Err(err).Msg("handler::UpdateProduct - Failed to parse request body")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error(err))
	}

	resp, err := h.service.UpdateProduct(c.Context(), &req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.Error(err))
	}

	return c.Status(fiber.StatusOK).JSON(response.Success(resp, ""))
}

func (h *productHandler) DeleteProduct(c *fiber.Ctx) error {
	var req entity.DeleteProductRequest
	req.Id = c.Params("id")

	err := h.service.DeleteProduct(c.Context(), &req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.Error(err))
	}

	return c.Status(fiber.StatusOK).JSON(response.Success(nil, ""))
}

func (h *productHandler) GetAllProducts(c *fiber.Ctx) error {
	var req entity.GetAllProductRequest
	if err := c.QueryParser(&req); err != nil {
		log.Warn().Err(err).Msg("handler::GetAllProducts - Failed to parse query params")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error(err))
	}

	resp, err := h.service.GetAllProducts(c.Context(), &req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.Error(err))
	}

	return c.Status(fiber.StatusOK).JSON(response.Success(resp, ""))
}
