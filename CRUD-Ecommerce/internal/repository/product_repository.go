package repository

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/ekas-7/CRUD-Ecommerce/internal/model"
	"github.com/google/uuid"
)

type ProductRepository interface {
	Create(product *model.Product) error
	GetByID(id uuid.UUID) (*model.Product, error)
	GetAll(params model.ProductQueryParams) ([]model.Product, error)
	GetByCategory(categoryID uuid.UUID) ([]model.Product, error)
	Update(product *model.Product) error
	UpdateStock(id uuid.UUID, stock int) error
	Delete(id uuid.UUID) error
}

type productRepository struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) ProductRepository {
	return &productRepository{db: db}
}

func (r *productRepository) Create(product *model.Product) error {
	query := `
		INSERT INTO products (id, name, description, price, stock, category_id, image_url, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id, created_at, updated_at
	`

	product.ID = uuid.New()
	product.CreatedAt = time.Now()
	product.UpdatedAt = time.Now()

	err := r.db.QueryRow(
		query,
		product.ID,
		product.Name,
		product.Description,
		product.Price,
		product.Stock,
		product.CategoryID,
		product.ImageURL,
		product.CreatedAt,
		product.UpdatedAt,
	).Scan(&product.ID, &product.CreatedAt, &product.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create product: %w", err)
	}

	return nil
}

func (r *productRepository) GetByID(id uuid.UUID) (*model.Product, error) {
	query := `
		SELECT p.id, p.name, p.description, p.price, p.stock, p.category_id, p.image_url, 
		       p.created_at, p.updated_at,
		       c.id, c.name, c.description, c.created_at, c.updated_at
		FROM products p
		LEFT JOIN categories c ON p.category_id = c.id
		WHERE p.id = $1
	`

	product := &model.Product{Category: &model.Category{}}
	var categoryID sql.NullString
	var categoryName sql.NullString
	var categoryDesc sql.NullString
	var categoryCreated sql.NullTime
	var categoryUpdated sql.NullTime

	err := r.db.QueryRow(query, id).Scan(
		&product.ID,
		&product.Name,
		&product.Description,
		&product.Price,
		&product.Stock,
		&product.CategoryID,
		&product.ImageURL,
		&product.CreatedAt,
		&product.UpdatedAt,
		&categoryID,
		&categoryName,
		&categoryDesc,
		&categoryCreated,
		&categoryUpdated,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("product not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get product: %w", err)
	}

	if categoryID.Valid {
		categoryUUID, _ := uuid.Parse(categoryID.String)
		product.Category.ID = categoryUUID
		product.Category.Name = categoryName.String
		product.Category.Description = categoryDesc.String
		product.Category.CreatedAt = categoryCreated.Time
		product.Category.UpdatedAt = categoryUpdated.Time
	} else {
		product.Category = nil
	}

	return product, nil
}

func (r *productRepository) GetAll(params model.ProductQueryParams) ([]model.Product, error) {
	query := `
		SELECT p.id, p.name, p.description, p.price, p.stock, p.category_id, p.image_url, 
		       p.created_at, p.updated_at
		FROM products p
		WHERE 1=1
	`

	var args []interface{}
	argPos := 1

	if params.CategoryID != uuid.Nil {
		query += fmt.Sprintf(" AND p.category_id = $%d", argPos)
		args = append(args, params.CategoryID)
		argPos++
	}

	if params.MinPrice > 0 {
		query += fmt.Sprintf(" AND p.price >= $%d", argPos)
		args = append(args, params.MinPrice)
		argPos++
	}

	if params.MaxPrice > 0 {
		query += fmt.Sprintf(" AND p.price <= $%d", argPos)
		args = append(args, params.MaxPrice)
		argPos++
	}

	if params.Search != "" {
		query += fmt.Sprintf(" AND (p.name ILIKE $%d OR p.description ILIKE $%d)", argPos, argPos)
		args = append(args, "%"+params.Search+"%")
		argPos++
	}

	query += " ORDER BY p.created_at DESC"

	if params.PageSize > 0 {
		query += fmt.Sprintf(" LIMIT $%d", argPos)
		args = append(args, params.PageSize)
		argPos++

		if params.Page > 0 {
			offset := (params.Page - 1) * params.PageSize
			query += fmt.Sprintf(" OFFSET $%d", argPos)
			args = append(args, offset)
		}
	}

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get products: %w", err)
	}
	defer rows.Close()

	var products []model.Product
	for rows.Next() {
		var product model.Product
		err := rows.Scan(
			&product.ID,
			&product.Name,
			&product.Description,
			&product.Price,
			&product.Stock,
			&product.CategoryID,
			&product.ImageURL,
			&product.CreatedAt,
			&product.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan product: %w", err)
		}
		products = append(products, product)
	}

	return products, nil
}

func (r *productRepository) GetByCategory(categoryID uuid.UUID) ([]model.Product, error) {
	return r.GetAll(model.ProductQueryParams{CategoryID: categoryID})
}

func (r *productRepository) Update(product *model.Product) error {
	updates := []string{}
	args := []interface{}{}
	argPos := 1

	if product.Name != "" {
		updates = append(updates, fmt.Sprintf("name = $%d", argPos))
		args = append(args, product.Name)
		argPos++
	}

	if product.Description != "" {
		updates = append(updates, fmt.Sprintf("description = $%d", argPos))
		args = append(args, product.Description)
		argPos++
	}

	if product.Price > 0 {
		updates = append(updates, fmt.Sprintf("price = $%d", argPos))
		args = append(args, product.Price)
		argPos++
	}

	if product.Stock >= 0 {
		updates = append(updates, fmt.Sprintf("stock = $%d", argPos))
		args = append(args, product.Stock)
		argPos++
	}

	if product.ImageURL != "" {
		updates = append(updates, fmt.Sprintf("image_url = $%d", argPos))
		args = append(args, product.ImageURL)
		argPos++
	}

	if len(updates) == 0 {
		return fmt.Errorf("no fields to update")
	}

	product.UpdatedAt = time.Now()
	updates = append(updates, fmt.Sprintf("updated_at = $%d", argPos))
	args = append(args, product.UpdatedAt)
	argPos++

	args = append(args, product.ID)
	query := fmt.Sprintf("UPDATE products SET %s WHERE id = $%d RETURNING updated_at",
		strings.Join(updates, ", "), argPos)

	err := r.db.QueryRow(query, args...).Scan(&product.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to update product: %w", err)
	}

	return nil
}

func (r *productRepository) UpdateStock(id uuid.UUID, stock int) error {
	query := `UPDATE products SET stock = $1, updated_at = $2 WHERE id = $3`

	result, err := r.db.Exec(query, stock, time.Now(), id)
	if err != nil {
		return fmt.Errorf("failed to update stock: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("product not found")
	}

	return nil
}

func (r *productRepository) Delete(id uuid.UUID) error {
	query := `DELETE FROM products WHERE id = $1`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete product: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("product not found")
	}

	return nil
}
