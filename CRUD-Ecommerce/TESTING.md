# API Testing with cURL

## Authentication

### Register
```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "password123",
    "first_name": "Test",
    "last_name": "User"
  }'
```

### Login
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "password123"
  }'
```

Save the token from the response for subsequent requests.

## Categories

### Get All Categories
```bash
curl http://localhost:8080/api/v1/categories
```

### Create Category (Admin)
```bash
curl -X POST http://localhost:8080/api/v1/categories \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{
    "name": "Electronics",
    "description": "Electronic devices and accessories"
  }'
```

## Products

### Get All Products
```bash
curl "http://localhost:8080/api/v1/products?page=1&page_size=10"
```

### Search Products
```bash
curl "http://localhost:8080/api/v1/products?search=laptop&min_price=500&max_price=2000"
```

### Create Product (Admin)
```bash
curl -X POST http://localhost:8080/api/v1/products \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{
    "name": "Gaming Laptop",
    "description": "High-performance gaming laptop",
    "price": 1299.99,
    "stock": 15,
    "category_id": "CATEGORY_UUID",
    "image_url": "https://example.com/laptop.jpg"
  }'
```

## Orders

### Create Order
```bash
curl -X POST http://localhost:8080/api/v1/orders \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{
    "items": [
      {
        "product_id": "PRODUCT_UUID",
        "quantity": 2
      }
    ]
  }'
```

### Get My Orders
```bash
curl http://localhost:8080/api/v1/orders \
  -H "Authorization: Bearer YOUR_TOKEN"
```

### Update Order Status
```bash
curl -X PUT http://localhost:8080/api/v1/orders/ORDER_UUID/status \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{
    "status": "processing"
  }'
```

## User Profile

### Get Profile
```bash
curl http://localhost:8080/api/v1/users/me \
  -H "Authorization: Bearer YOUR_TOKEN"
```

### Update Profile
```bash
curl -X PUT http://localhost:8080/api/v1/users/me \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{
    "first_name": "Updated",
    "last_name": "Name"
  }'
```
