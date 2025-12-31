# E-Commerce CRUD API - Go

A modern, scalable, and modular RESTful API for an e-commerce platform built with Go, Gin framework, and PostgreSQL.

## 🚀 Features

- **User Authentication & Authorization** - JWT-based authentication with role-based access control
- **Product Management** - Complete CRUD operations for products with category support
- **Order Management** - Create, track, and manage customer orders
- **Category Management** - Organize products with categories
- **Stock Management** - Automatic inventory tracking
- **Clean Architecture** - Follows best practices with clear separation of concerns
- **Scalable Design** - Modular structure ready for growth

## 📁 Project Structure

```
CRUD-Ecommerce/
├── cmd/
│   └── api/
│       └── main.go              # Application entry point
├── internal/
│   ├── config/
│   │   └── config.go            # Configuration management
│   ├── database/
│   │   └── postgres.go          # Database connection & migrations
│   ├── handler/
│   │   ├── user_handler.go      # User HTTP handlers
│   │   ├── product_handler.go   # Product HTTP handlers
│   │   ├── category_handler.go  # Category HTTP handlers
│   │   ├── order_handler.go     # Order HTTP handlers
│   │   └── handler.go           # Handler utilities
│   ├── middleware/
│   │   ├── auth.go              # Authentication middleware
│   │   ├── cors.go              # CORS middleware
│   │   ├── logger.go            # Logging middleware
│   │   └── recovery.go          # Panic recovery middleware
│   ├── model/
│   │   ├── user.go              # User domain models
│   │   ├── product.go           # Product domain models
│   │   ├── category.go          # Category domain models
│   │   └── order.go             # Order domain models
│   ├── repository/
│   │   ├── user_repository.go   # User data access
│   │   ├── product_repository.go # Product data access
│   │   ├── category_repository.go # Category data access
│   │   ├── order_repository.go  # Order data access
│   │   └── repository.go        # Repository container
│   └── service/
│       ├── user_service.go      # User business logic
│       ├── product_service.go   # Product business logic
│       ├── category_service.go  # Category business logic
│       ├── order_service.go     # Order business logic
│       └── service.go           # Service container
├── migrations/
│   ├── 001_initial_schema.sql   # Database schema
│   └── 002_seed_data.sql        # Sample data
├── .env.example                 # Environment variables template
├── .gitignore                   # Git ignore rules
├── go.mod                       # Go module definition
└── README.md                    # This file
```

## 🛠️ Tech Stack

- **Language**: Go 1.21+
- **Web Framework**: Gin
- **Database**: PostgreSQL
- **Authentication**: JWT (golang-jwt/jwt)
- **Password Hashing**: bcrypt
- **Environment**: godotenv
- **Validation**: go-playground/validator

## 📋 Prerequisites

- Go 1.21 or higher
- PostgreSQL 12 or higher
- Git

## 🔧 Installation

### Option 1: Using Docker (Recommended - Easiest)

1. **Clone the repository**
```bash
git clone https://github.com/ekas-7/CRUD-Ecommerce.git
cd CRUD-Ecommerce
```

2. **Install Go dependencies**
```bash
go mod download
```

3. **Start PostgreSQL with Docker**
```bash
docker-compose up -d postgres
# Or using Makefile
make db-start
```

4. **Wait for database to be ready (3-5 seconds)**
```bash
sleep 3
```

5. **Run the application**
```bash
go run cmd/api/main.go
# Or using Makefile
make run
```

The server will start on `http://localhost:8080`

### Option 2: Local PostgreSQL Installation

1. **Clone the repository**
```bash
git clone https://github.com/ekas-7/CRUD-Ecommerce.git
cd CRUD-Ecommerce
```

2. **Install dependencies**
```bash
go mod download
```

3. **Set up PostgreSQL database**
```bash
# Install PostgreSQL (if not already installed)
# macOS
brew install postgresql@15
brew services start postgresql@15

# Create database
createdb ecommerce_db
```

4. **Configure environment variables**
```bash
cp .env.example .env
# Edit .env - change DB_PORT to 5432 if using local PostgreSQL
```

5. **Run the application**
```bash
go run cmd/api/main.go
```

## 🌍 Environment Variables

### Using Docker (Default Configuration)

The `.env` file is already configured for Docker:

```env
# Server Configuration
SERVER_PORT=8080
SERVER_HOST=localhost

# Database Configuration (Docker)
DB_HOST=localhost
DB_PORT=5433
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=ecommerce_db
DB_SSLMODE=disable

# JWT Configuration
JWT_SECRET=your-super-secret-key-change-this-in-production
JWT_EXPIRY=24h

# Application
APP_ENV=development
```

### Using Local PostgreSQL

If you're using a local PostgreSQL installation instead of Docker, change:
```env
DB_PORT=5432  # Change from 5433 to 5432
```

## 📚 API Documentation

### Base URL
```
http://localhost:8080/api/v1
```

### Authentication

#### Register User
```http
POST /api/v1/auth/register
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "password123",
  "first_name": "John",
  "last_name": "Doe"
}
```

#### Login
```http
POST /api/v1/auth/login
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "password123"
}

Response:
{
  "token": "eyJhbGciOiJIUzI1NiIs...",
  "user": {
    "id": "uuid",
    "email": "user@example.com",
    "first_name": "John",
    "last_name": "Doe",
    "role": "user"
  }
}
```

### Categories

#### Get All Categories
```http
GET /api/v1/categories
```

#### Get Category by ID
```http
GET /api/v1/categories/:id
```

#### Create Category (Admin Only)
```http
POST /api/v1/categories
Authorization: Bearer <token>
Content-Type: application/json

{
  "name": "Electronics",
  "description": "Electronic devices"
}
```

#### Update Category (Admin Only)
```http
PUT /api/v1/categories/:id
Authorization: Bearer <token>
Content-Type: application/json

{
  "name": "Updated Electronics",
  "description": "Updated description"
}
```

#### Delete Category (Admin Only)
```http
DELETE /api/v1/categories/:id
Authorization: Bearer <token>
```

### Products

#### Get All Products
```http
GET /api/v1/products?page=1&page_size=10&search=laptop&min_price=100&max_price=1000
```

#### Get Product by ID
```http
GET /api/v1/products/:id
```

#### Get Products by Category
```http
GET /api/v1/products/category/:categoryId
```

#### Create Product (Admin Only)
```http
POST /api/v1/products
Authorization: Bearer <token>
Content-Type: application/json

{
  "name": "Laptop",
  "description": "High-performance laptop",
  "price": 999.99,
  "stock": 10,
  "category_id": "category-uuid",
  "image_url": "https://example.com/image.jpg"
}
```

#### Update Product (Admin Only)
```http
PUT /api/v1/products/:id
Authorization: Bearer <token>
Content-Type: application/json

{
  "name": "Updated Laptop",
  "price": 899.99,
  "stock": 15
}
```

#### Delete Product (Admin Only)
```http
DELETE /api/v1/products/:id
Authorization: Bearer <token>
```

### Orders

#### Create Order
```http
POST /api/v1/orders
Authorization: Bearer <token>
Content-Type: application/json

{
  "items": [
    {
      "product_id": "product-uuid",
      "quantity": 2
    }
  ]
}
```

#### Get User Orders
```http
GET /api/v1/orders
Authorization: Bearer <token>
```

#### Get Order by ID
```http
GET /api/v1/orders/:id
Authorization: Bearer <token>
```

#### Update Order Status
```http
PUT /api/v1/orders/:id/status
Authorization: Bearer <token>
Content-Type: application/json

{
  "status": "processing"
}

# Valid statuses: pending, processing, shipped, delivered, cancelled
```

#### Cancel Order
```http
DELETE /api/v1/orders/:id
Authorization: Bearer <token>
```

#### Get All Orders (Admin Only)
```http
GET /api/v1/orders/all
Authorization: Bearer <token>
```

### User Profile

#### Get Profile
```http
GET /api/v1/users/me
Authorization: Bearer <token>
```

#### Update Profile
```http
PUT /api/v1/users/me
Authorization: Bearer <token>
Content-Type: application/json

{
  "first_name": "Jane",
  "last_name": "Smith"
}
```

#### Delete Account
```http
DELETE /api/v1/users/me
Authorization: Bearer <token>
```

## 🏗️ Architecture

This project follows **Clean Architecture** principles:

1. **Handler Layer** - HTTP request/response handling
2. **Service Layer** - Business logic and validation
3. **Repository Layer** - Data access and persistence
4. **Model Layer** - Domain entities and DTOs

### Key Design Patterns

- **Dependency Injection** - Services and repositories are injected
- **Repository Pattern** - Abstract data access
- **Middleware Chain** - Authentication, logging, CORS, recovery
- **Interface-based Design** - Easy testing and mocking

## 🔐 Security Features

- **JWT Authentication** - Secure token-based auth
- **Password Hashing** - bcrypt for password security
- **Role-based Access Control** - Admin and user roles
- **CORS Protection** - Configurable cross-origin requests
- **SQL Injection Prevention** - Parameterized queries
- **Request Validation** - Input validation at handler level

## 🧪 Testing

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests for specific package
go test ./internal/service/...
```

## 🚢 Deployment

### Docker (Optional)

Create a `Dockerfile`:
```dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY go.* ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o main cmd/api/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/main .
COPY --from=builder /app/.env .
EXPOSE 8080
CMD ["./main"]
```

Build and run:
```bash
docker build -t ecommerce-api .
docker run -p 8080:8080 ecommerce-api
```

## 📈 Performance Tips

- Use connection pooling (configured in database setup)
- Implement caching for frequently accessed data
- Use database indexes (already configured)
- Implement pagination for large datasets
- Use prepared statements for repeated queries

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📝 License

This project is licensed under the MIT License.

## 👤 Author

**ekas-7**
- GitHub: [@ekas-7](https://github.com/ekas-7)

## 🙏 Acknowledgments

- Gin Web Framework
- PostgreSQL
- The Go Community

## 📞 Support

For support, email support@example.com or create an issue in the repository.

---

**Happy Coding! 🚀**
