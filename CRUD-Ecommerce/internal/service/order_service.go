package service

import (
	"fmt"

	"github.com/ekas-7/CRUD-Ecommerce/internal/model"
	"github.com/ekas-7/CRUD-Ecommerce/internal/repository"
	"github.com/google/uuid"
)

type OrderService interface {
	Create(userID uuid.UUID, req *model.OrderCreateRequest) (*model.Order, error)
	GetByID(orderID, userID uuid.UUID, isAdmin bool) (*model.Order, error)
	GetUserOrders(userID uuid.UUID) ([]model.Order, error)
	GetAllOrders() ([]model.Order, error)
	UpdateStatus(orderID uuid.UUID, status model.OrderStatus) (*model.Order, error)
	Cancel(orderID, userID uuid.UUID, isAdmin bool) error
}

type orderService struct {
	orderRepo   repository.OrderRepository
	productRepo repository.ProductRepository
}

func NewOrderService(orderRepo repository.OrderRepository, productRepo repository.ProductRepository) OrderService {
	return &orderService{
		orderRepo:   orderRepo,
		productRepo: productRepo,
	}
}

func (s *orderService) Create(userID uuid.UUID, req *model.OrderCreateRequest) (*model.Order, error) {
	if len(req.Items) == 0 {
		return nil, fmt.Errorf("order must contain at least one item")
	}

	order := &model.Order{
		UserID: userID,
		Status: model.OrderStatusPending,
		Items:  []model.OrderItem{},
	}

	var totalPrice float64

	// Process each item
	for _, itemReq := range req.Items {
		// Get product details
		product, err := s.productRepo.GetByID(itemReq.ProductID)
		if err != nil {
			return nil, fmt.Errorf("product %s not found: %w", itemReq.ProductID, err)
		}

		// Check stock availability
		if product.Stock < itemReq.Quantity {
			return nil, fmt.Errorf("insufficient stock for product %s. Available: %d, Requested: %d",
				product.Name, product.Stock, itemReq.Quantity)
		}

		// Calculate item price
		itemPrice := product.Price * float64(itemReq.Quantity)
		totalPrice += itemPrice

		// Create order item
		orderItem := model.OrderItem{
			ProductID: itemReq.ProductID,
			Quantity:  itemReq.Quantity,
			Price:     product.Price,
		}

		order.Items = append(order.Items, orderItem)

		// Update product stock
		newStock := product.Stock - itemReq.Quantity
		if err := s.productRepo.UpdateStock(product.ID, newStock); err != nil {
			return nil, fmt.Errorf("failed to update stock: %w", err)
		}
	}

	order.TotalPrice = totalPrice

	// Create order in database
	if err := s.orderRepo.Create(order); err != nil {
		return nil, fmt.Errorf("failed to create order: %w", err)
	}

	// Fetch full order with product details
	return s.orderRepo.GetByID(order.ID)
}

func (s *orderService) GetByID(orderID, userID uuid.UUID, isAdmin bool) (*model.Order, error) {
	order, err := s.orderRepo.GetByID(orderID)
	if err != nil {
		return nil, err
	}

	// Check if user owns the order or is admin
	if !isAdmin && order.UserID != userID {
		return nil, fmt.Errorf("access denied: order does not belong to user")
	}

	return order, nil
}

func (s *orderService) GetUserOrders(userID uuid.UUID) ([]model.Order, error) {
	return s.orderRepo.GetByUserID(userID)
}

func (s *orderService) GetAllOrders() ([]model.Order, error) {
	return s.orderRepo.GetAll()
}

func (s *orderService) UpdateStatus(orderID uuid.UUID, status model.OrderStatus) (*model.Order, error) {
	// Validate order exists
	order, err := s.orderRepo.GetByID(orderID)
	if err != nil {
		return nil, err
	}

	// Validate status transition
	if order.Status == model.OrderStatusCancelled {
		return nil, fmt.Errorf("cannot update cancelled order")
	}

	if order.Status == model.OrderStatusDelivered && status != model.OrderStatusDelivered {
		return nil, fmt.Errorf("cannot change status of delivered order")
	}

	if err := s.orderRepo.UpdateStatus(orderID, status); err != nil {
		return nil, fmt.Errorf("failed to update order status: %w", err)
	}

	return s.orderRepo.GetByID(orderID)
}

func (s *orderService) Cancel(orderID, userID uuid.UUID, isAdmin bool) error {
	order, err := s.orderRepo.GetByID(orderID)
	if err != nil {
		return err
	}

	// Check if user owns the order or is admin
	if !isAdmin && order.UserID != userID {
		return fmt.Errorf("access denied: order does not belong to user")
	}

	// Only pending or processing orders can be cancelled
	if order.Status != model.OrderStatusPending && order.Status != model.OrderStatusProcessing {
		return fmt.Errorf("order cannot be cancelled in current status: %s", order.Status)
	}

	// Restore product stock
	for _, item := range order.Items {
		product, err := s.productRepo.GetByID(item.ProductID)
		if err != nil {
			continue // Skip if product not found
		}

		newStock := product.Stock + item.Quantity
		s.productRepo.UpdateStock(product.ID, newStock)
	}

	// Update order status to cancelled
	return s.orderRepo.UpdateStatus(orderID, model.OrderStatusCancelled)
}
