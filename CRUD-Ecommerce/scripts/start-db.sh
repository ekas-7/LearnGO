#!/bin/bash

echo "🚀 Starting PostgreSQL with Docker..."

# Stop any existing container
docker-compose down 2>/dev/null

# Start PostgreSQL
docker-compose up -d postgres

# Wait for PostgreSQL to be ready
echo "⏳ Waiting for PostgreSQL to be ready..."
sleep 3

# Check if database is ready
until docker exec ecommerce_postgres pg_isready -U postgres -d ecommerce_db > /dev/null 2>&1; do
  echo "Waiting for database..."
  sleep 1
done

echo "✅ PostgreSQL is ready!"
echo ""
echo "📊 Database Connection Details:"
echo "  Host: localhost"
echo "  Port: 5432"
echo "  Database: ecommerce_db"
echo "  User: postgres"
echo "  Password: postgres"
echo ""
echo "🔗 Connection string: postgresql://postgres:postgres@localhost:5432/ecommerce_db"
echo ""
echo "📝 To view logs: docker logs ecommerce_postgres -f"
echo "🛑 To stop: docker-compose down"
echo ""
echo "✨ Now you can run: go run cmd/api/main.go"
