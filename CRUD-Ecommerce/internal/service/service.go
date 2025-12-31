package service

type Services struct {
	User     UserService
	Product  ProductService
	Category CategoryService
	Order    OrderService
}
