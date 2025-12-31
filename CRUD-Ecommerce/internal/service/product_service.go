package service

import (
	"fmt"

	"github.com/ekas-7/CRUD-Ecommerce/internal/model"
	"github.com/ekas-7/CRUD-Ecommerce/internal/repository"
	"github.com/google/uuid"
)

type ProductService interface {
	Create(req *model.ProductCreateRequest) (*model.Product, error)
	GetByID(id uuid.UUID) (*model.Product, error)
	GetAll(params model.ProductQueryParams) ([]model.Product, error)
	GetByCategory(categoryID uuid.UUID) ([]model.Product, error)
	Update(id uuid.UUID, req *model.ProductUpdateRequest) (*model.Product, error)
	Delete(id uuid.UUID) error
}

type productService struct {
	repo repository.ProductRepository
}

func NewProductService(repo repository.ProductRepository) ProductService {
	return &productService{repo: repo}
}

func (s *productService) Create(req *model.ProductCreateRequest) (*model.Product, error) {
	product := &model.Product{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Stock:       req.Stock,
		CategoryID:  req.CategoryID,
		ImageURL:    req.ImageURL,
	}

	if err := s.repo.Create(product); err != nil {
		return nil, fmt.Errorf("failed to create product: %w", err)
	}

	return product, nil
}

func (s *productService) GetByID(id uuid.UUID) (*model.Product, error) {
	return s.repo.GetByID(id)
}

func (s *productService) GetAll(params model.ProductQueryParams) ([]model.Product, error) {
	return s.repo.GetAll(params)
}

func (s *productService) GetByCategory(categoryID uuid.UUID) ([]model.Product, error) {
	return s.repo.GetByCategory(categoryID)
}

func (s *productService) Update(id uuid.UUID, req *model.ProductUpdateRequest) (*model.Product, error) {
	product, err := s.repo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("product not found: %w", err)
	}

	if req.Name != "" {
		product.Name = req.Name
	}
	if req.Description != "" {
		product.Description = req.Description
	}
	if req.Price > 0 {
		product.Price = req.Price
	}
	if req.Stock >= 0 {
		product.Stock = req.Stock
	}
	if req.ImageURL != "" {
		product.ImageURL = req.ImageURL
	}

	if err := s.repo.Update(product); err != nil {
		return nil, fmt.Errorf("failed to update product: %w", err)
	}

	return product, nil
}

func (s *productService) Delete(id uuid.UUID) error {
	return s.repo.Delete(id)
}
