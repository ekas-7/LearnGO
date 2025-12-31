package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/ekas-7/CRUD-Ecommerce/internal/model"
	"github.com/google/uuid"
)

type OrderRepository interface {
	Create(order *model.Order) error
	GetByID(id uuid.UUID) (*model.Order, error)
	GetByUserID(userID uuid.UUID) ([]model.Order, error)
	GetAll() ([]model.Order, error)
	UpdateStatus(id uuid.UUID, status model.OrderStatus) error
	Delete(id uuid.UUID) error
}

type orderRepository struct {
	db *sql.DB
}

func NewOrderRepository(db *sql.DB) OrderRepository {
	return &orderRepository{db: db}
}

func (r *orderRepository) Create(order *model.Order) error {
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Insert order
	orderQuery := `
		INSERT INTO orders (id, user_id, status, total_price, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at, updated_at
	`

	order.ID = uuid.New()
	order.CreatedAt = time.Now()
	order.UpdatedAt = time.Now()

	if order.Status == "" {
		order.Status = model.OrderStatusPending
	}

	err = tx.QueryRow(
		orderQuery,
		order.ID,
		order.UserID,
		order.Status,
		order.TotalPrice,
		order.CreatedAt,
		order.UpdatedAt,
	).Scan(&order.ID, &order.CreatedAt, &order.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create order: %w", err)
	}

	// Insert order items
	itemQuery := `
		INSERT INTO order_items (id, order_id, product_id, quantity, price, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at
	`

	for i := range order.Items {
		item := &order.Items[i]
		item.ID = uuid.New()
		item.OrderID = order.ID
		item.CreatedAt = time.Now()

		err = tx.QueryRow(
			itemQuery,
			item.ID,
			item.OrderID,
			item.ProductID,
			item.Quantity,
			item.Price,
			item.CreatedAt,
		).Scan(&item.ID, &item.CreatedAt)

		if err != nil {
			return fmt.Errorf("failed to create order item: %w", err)
		}
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (r *orderRepository) GetByID(id uuid.UUID) (*model.Order, error) {
	orderQuery := `
		SELECT id, user_id, status, total_price, created_at, updated_at
		FROM orders
		WHERE id = $1
	`

	order := &model.Order{}
	err := r.db.QueryRow(orderQuery, id).Scan(
		&order.ID,
		&order.UserID,
		&order.Status,
		&order.TotalPrice,
		&order.CreatedAt,
		&order.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("order not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get order: %w", err)
	}

	// Get order items
	itemsQuery := `
		SELECT oi.id, oi.order_id, oi.product_id, oi.quantity, oi.price, oi.created_at,
		       p.id, p.name, p.description, p.price, p.stock, p.category_id, p.image_url, p.created_at, p.updated_at
		FROM order_items oi
		LEFT JOIN products p ON oi.product_id = p.id
		WHERE oi.order_id = $1
	`

	rows, err := r.db.Query(itemsQuery, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get order items: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var item model.OrderItem
		item.Product = &model.Product{}

		err := rows.Scan(
			&item.ID,
			&item.OrderID,
			&item.ProductID,
			&item.Quantity,
			&item.Price,
			&item.CreatedAt,
			&item.Product.ID,
			&item.Product.Name,
			&item.Product.Description,
			&item.Product.Price,
			&item.Product.Stock,
			&item.Product.CategoryID,
			&item.Product.ImageURL,
			&item.Product.CreatedAt,
			&item.Product.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan order item: %w", err)
		}

		order.Items = append(order.Items, item)
	}

	return order, nil
}

func (r *orderRepository) GetByUserID(userID uuid.UUID) ([]model.Order, error) {
	query := `
		SELECT id, user_id, status, total_price, created_at, updated_at
		FROM orders
		WHERE user_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get orders: %w", err)
	}
	defer rows.Close()

	var orders []model.Order
	for rows.Next() {
		var order model.Order
		err := rows.Scan(
			&order.ID,
			&order.UserID,
			&order.Status,
			&order.TotalPrice,
			&order.CreatedAt,
			&order.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan order: %w", err)
		}

		// Get items for each order
		fullOrder, err := r.GetByID(order.ID)
		if err != nil {
			return nil, err
		}
		order.Items = fullOrder.Items

		orders = append(orders, order)
	}

	return orders, nil
}

func (r *orderRepository) GetAll() ([]model.Order, error) {
	query := `
		SELECT id, user_id, status, total_price, created_at, updated_at
		FROM orders
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get orders: %w", err)
	}
	defer rows.Close()

	var orders []model.Order
	for rows.Next() {
		var order model.Order
		err := rows.Scan(
			&order.ID,
			&order.UserID,
			&order.Status,
			&order.TotalPrice,
			&order.CreatedAt,
			&order.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan order: %w", err)
		}

		// Get items for each order
		fullOrder, err := r.GetByID(order.ID)
		if err != nil {
			return nil, err
		}
		order.Items = fullOrder.Items

		orders = append(orders, order)
	}

	return orders, nil
}

func (r *orderRepository) UpdateStatus(id uuid.UUID, status model.OrderStatus) error {
	query := `
		UPDATE orders
		SET status = $1, updated_at = $2
		WHERE id = $3
		RETURNING updated_at
	`

	var updatedAt time.Time
	err := r.db.QueryRow(query, status, time.Now(), id).Scan(&updatedAt)
	if err != nil {
		return fmt.Errorf("failed to update order status: %w", err)
	}

	return nil
}

func (r *orderRepository) Delete(id uuid.UUID) error {
	query := `DELETE FROM orders WHERE id = $1`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete order: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("order not found")
	}

	return nil
}
