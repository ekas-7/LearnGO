package repository

import "database/sql"

type Repositories struct {
	User     UserRepository
	Product  ProductRepository
	Category CategoryRepository
	Order    OrderRepository
}

func NewRepositories(db *sql.DB) *Repositories {
	return &Repositories{
		User:     NewUserRepository(db),
		Product:  NewProductRepository(db),
		Category: NewCategoryRepository(db),
		Order:    NewOrderRepository(db),
	}
}
