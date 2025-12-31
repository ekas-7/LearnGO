# Quick Start Guide

## Option 1: Docker Setup (Recommended - Easiest) ⭐

### Prerequisites
- Go 1.21+
- Docker & Docker Compose

### Steps

1. **Start PostgreSQL with Docker**
   ```bash
   cd /Users/ekas/Desktop/LearnGO/CRUD-Ecommerce
   docker-compose up -d postgres
   ```
   
   Or using Makefile:
   ```bash
   make db-start
   ```

2. **Wait for Database to be Ready (3-5 seconds)**
   ```bash
   sleep 3
   ```
   
   Verify it's ready:
   ```bash
   docker exec ecommerce_postgres pg_isready -U postgres -d ecommerce_db
   ```

3. **Install Go Dependencies** (first time only)
   ```bash
   go mod download
   ```

4. **Run the Application**
   ```bash
   go run cmd/api/main.go
   ```
   
   Or using Makefile:
   ```bash
   make run
   ```

5. **Test the API**
   ```bash
   # Health check
   curl http://localhost:8080/health
   
   # Register a user
   curl -X POST http://localhost:8080/api/v1/auth/register \
     -H "Content-Type: application/json" \
     -d '{
       "email": "test@example.com",
       "password": "password123",
       "first_name": "Test",
       "last_name": "User"
     }'
   ```

### Database Management Commands

```bash
# Start database
make db-start

# Stop database
make db-stop

# View logs
make db-logs

# Access PostgreSQL shell
make db-shell

# Reset database (delete all data)
make db-reset
```

### Important Notes

- **Database runs on port 5433** (not 5432) to avoid conflicts
- **Connection string:** `postgresql://postgres:postgres@localhost:5433/ecommerce_db`
- The `.env` file is already configured for Docker

---

## Option 2: Local PostgreSQL Setup

### Prerequisites
- Go 1.21+
- PostgreSQL 12+

### Steps

1. **Install PostgreSQL** (if not already installed)
   ```bash
   # macOS
   brew install postgresql@15
   brew services start postgresql@15
   ```

2. **Create Database**
   ```bash
   createdb ecommerce_db
   ```

3. **Install Dependencies**
   ```bash
   cd /Users/ekas/Desktop/LearnGO/CRUD-Ecommerce
   go mod download
   ```

4. **Configure Environment**
   - The `.env` file is already created with default values
   - Update `DB_PASSWORD` if your PostgreSQL has a different password
   - Update `JWT_SECRET` for production

5. **Run the Application**
   ```bash
   go run cmd/api/main.go
   ```
   
   Or using Makefile:
   ```bash
   make run
   ```

6. **Test the API**
   ```bash
   # Health check
   curl http://localhost:8080/health
   
   # Register a user
   curl -X POST http://localhost:8080/api/v1/auth/register \
     -H "Content-Type: application/json" \
     -d '{
       "email": "test@example.com",
       "password": "password123",
       "first_name": "Test",
       "last_name": "User"
     }'
   ```

## Option 2: Docker (Recommended for Production)

### Prerequisites
- Docker
- Docker Compose

### Steps

1. **Start Services**
   ```bash
   docker-compose up -d
   ```
   This will:
   - Start PostgreSQL database
   - Run database migrations
   - Start the API server

2. **Check Logs**
   ```bash
   docker-compose logs -f api
   ```

3. **Stop Services**
   ```bash
   docker-compose down
   ```

## Option 3: Hot Reload Development

### Prerequisites
- Air (for hot reload)

### Steps

1. **Install Air**
   ```bash
   go install github.com/cosmtrek/air@latest
   ```
   
   Or using Makefile:
   ```bash
   make install-air
   ```

2. **Run with Hot Reload**
   ```bash
   air
   ```
   
   Or using Makefile:
   ```bash
   make dev
   ```

## Useful Commands

```bash
# Build the application
make build

# Run tests
make test

# Format code
make fmt

# Run linters
make lint

# Clean build artifacts
make clean

# View all available commands
make help
```

## Next Steps

1. **Read the API Documentation** - Check `README.md` for complete API reference
2. **Test the Endpoints** - Use the examples in `TESTING.md`
3. **Customize** - Modify the code to fit your specific requirements
4. **Deploy** - Use Docker for production deployment

## Troubleshooting

### Database Connection Issues
- Ensure PostgreSQL is running: `pg_isready`
- Check connection settings in `.env`
- Verify database exists: `psql -l | grep ecommerce_db`

### Port Already in Use
- Change `SERVER_PORT` in `.env` to a different port
- Or kill the process using port 8080: `lsof -ti:8080 | xargs kill`

### Module Issues
```bash
go mod tidy
go mod download
```

## Default Access

- **API URL**: http://localhost:8080
- **Health Check**: http://localhost:8080/health
- **Database**: localhost:5432/ecommerce_db

## Creating an Admin User

After starting the app, register a user and manually update their role in the database:

```sql
psql ecommerce_db
UPDATE users SET role = 'admin' WHERE email = 'your-email@example.com';
```

Or create directly:
```bash
# First, hash your password using bcrypt (you can use an online tool)
# Then insert:
psql ecommerce_db -c "INSERT INTO users (email, password, first_name, last_name, role) VALUES ('admin@example.com', 'HASHED_PASSWORD', 'Admin', 'User', 'admin');"
```

Happy coding! 🚀
