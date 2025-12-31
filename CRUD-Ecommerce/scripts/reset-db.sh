#!/bin/bash

echo "🔄 Restarting PostgreSQL and cleaning data..."

# Stop and remove containers and volumes
docker-compose down -v

# Start fresh
docker-compose up -d postgres

# Wait for PostgreSQL
echo "⏳ Waiting for PostgreSQL to be ready..."
sleep 3

until docker exec ecommerce_postgres pg_isready -U postgres -d ecommerce_db > /dev/null 2>&1; do
  echo "Waiting for database..."
  sleep 1
done

echo "✅ PostgreSQL reset complete!"
echo "✨ Database is fresh and ready to use"
