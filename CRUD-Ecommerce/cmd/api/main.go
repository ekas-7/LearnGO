package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/ekas-7/CRUD-Ecommerce/internal/config"
	"github.com/ekas-7/CRUD-Ecommerce/internal/database"
	"github.com/ekas-7/CRUD-Ecommerce/internal/handler"
	"github.com/ekas-7/CRUD-Ecommerce/internal/middleware"
	"github.com/ekas-7/CRUD-Ecommerce/internal/repository"
	"github.com/ekas-7/CRUD-Ecommerce/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found, using system environment variables")
	}

	// Load configuration
	cfg := config.Load()

	// Initialize database
	db, err := database.NewPostgresDB(cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Run migrations
	if err := database.RunMigrations(db); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Initialize repositories
	repos := initRepositories(db)

	// Initialize services
	services := initServices(repos, cfg)

	// Initialize handlers
	handlers := initHandlers(services)

	// Setup router
	router := setupRouter(handlers, cfg)

	// Start server
	addr := fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port)
	log.Printf("Server starting on %s", addr)
	if err := router.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func initRepositories(db *sql.DB) *repository.Repositories {
	return &repository.Repositories{
		User:     repository.NewUserRepository(db),
		Product:  repository.NewProductRepository(db),
		Category: repository.NewCategoryRepository(db),
		Order:    repository.NewOrderRepository(db),
	}
}

func initServices(repos *repository.Repositories, cfg *config.Config) *service.Services {
	return &service.Services{
		User:     service.NewUserService(repos.User, cfg.JWT.Secret, cfg.JWT.Expiry),
		Product:  service.NewProductService(repos.Product),
		Category: service.NewCategoryService(repos.Category),
		Order:    service.NewOrderService(repos.Order, repos.Product),
	}
}

func initHandlers(services *service.Services) *handler.Handlers {
	return &handler.Handlers{
		User:     handler.NewUserHandler(services.User),
		Product:  handler.NewProductHandler(services.Product),
		Category: handler.NewCategoryHandler(services.Category),
		Order:    handler.NewOrderHandler(services.Order),
	}
}

func setupRouter(handlers *handler.Handlers, cfg *config.Config) *gin.Engine {
	if cfg.App.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()

	// Apply global middleware
	router.Use(middleware.CORS())
	router.Use(middleware.Logger())
	router.Use(middleware.Recovery())

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Public routes
		auth := v1.Group("/auth")
		{
			auth.POST("/register", handlers.User.Register)
			auth.POST("/login", handlers.User.Login)
		}

		// Categories (public read)
		categories := v1.Group("/categories")
		{
			categories.GET("", handlers.Category.GetAll)
			categories.GET("/:id", handlers.Category.GetByID)
		}

		// Products (public read)
		products := v1.Group("/products")
		{
			products.GET("", handlers.Product.GetAll)
			products.GET("/:id", handlers.Product.GetByID)
			products.GET("/category/:categoryId", handlers.Product.GetByCategory)
		}

		// Protected routes
		protected := v1.Group("")
		protected.Use(middleware.AuthMiddleware(cfg.JWT.Secret))
		{
			// User routes
			users := protected.Group("/users")
			{
				users.GET("/me", handlers.User.GetProfile)
				users.PUT("/me", handlers.User.UpdateProfile)
				users.DELETE("/me", handlers.User.DeleteAccount)
			}

			// Admin product routes
			adminProducts := protected.Group("/products")
			adminProducts.Use(middleware.AdminMiddleware())
			{
				adminProducts.POST("", handlers.Product.Create)
				adminProducts.PUT("/:id", handlers.Product.Update)
				adminProducts.DELETE("/:id", handlers.Product.Delete)
			}

			// Admin category routes
			adminCategories := protected.Group("/categories")
			adminCategories.Use(middleware.AdminMiddleware())
			{
				adminCategories.POST("", handlers.Category.Create)
				adminCategories.PUT("/:id", handlers.Category.Update)
				adminCategories.DELETE("/:id", handlers.Category.Delete)
			}

			// Order routes
			orders := protected.Group("/orders")
			{
				orders.POST("", handlers.Order.Create)
				orders.GET("", handlers.Order.GetUserOrders)
				orders.GET("/:id", handlers.Order.GetByID)
				orders.PUT("/:id/status", handlers.Order.UpdateStatus)
				orders.DELETE("/:id", handlers.Order.Cancel)
			}

			// Admin order routes
			adminOrders := protected.Group("/orders")
			adminOrders.Use(middleware.AdminMiddleware())
			{
				adminOrders.GET("/all", handlers.Order.GetAllOrders)
			}
		}
	}

	return router
}
