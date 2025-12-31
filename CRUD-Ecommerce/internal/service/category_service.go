package service

import (
	"fmt"

	"github.com/ekas-7/CRUD-Ecommerce/internal/model"
	"github.com/ekas-7/CRUD-Ecommerce/internal/repository"
	"github.com/google/uuid"
)

type CategoryService interface {
	Create(req *model.CategoryCreateRequest) (*model.Category, error)
	GetByID(id uuid.UUID) (*model.Category, error)
	GetAll() ([]model.Category, error)
	Update(id uuid.UUID, req *model.CategoryUpdateRequest) (*model.Category, error)
	Delete(id uuid.UUID) error
}

type categoryService struct {
	repo repository.CategoryRepository
}

func NewCategoryService(repo repository.CategoryRepository) CategoryService {
	return &categoryService{repo: repo}
}

func (s *categoryService) Create(req *model.CategoryCreateRequest) (*model.Category, error) {
	category := &model.Category{
		Name:        req.Name,
		Description: req.Description,
	}

	if err := s.repo.Create(category); err != nil {
		return nil, fmt.Errorf("failed to create category: %w", err)
	}

	return category, nil
}

func (s *categoryService) GetByID(id uuid.UUID) (*model.Category, error) {
	return s.repo.GetByID(id)
}

func (s *categoryService) GetAll() ([]model.Category, error) {
	return s.repo.GetAll()
}

func (s *categoryService) Update(id uuid.UUID, req *model.CategoryUpdateRequest) (*model.Category, error) {
	category, err := s.repo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("category not found: %w", err)
	}

	if req.Name != "" {
		category.Name = req.Name
	}
	if req.Description != "" {
		category.Description = req.Description
	}

	if err := s.repo.Update(category); err != nil {
		return nil, fmt.Errorf("failed to update category: %w", err)
	}

	return category, nil
}

func (s *categoryService) Delete(id uuid.UUID) error {
	return s.repo.Delete(id)
}
