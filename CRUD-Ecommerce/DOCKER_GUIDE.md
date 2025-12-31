# Docker Setup Guide for E-Commerce API

## 🐳 Quick Start with Docker

### Prerequisites
- Docker installed ([Get Docker](https://docs.docker.com/get-docker/))
- Docker Compose installed (usually comes with Docker Desktop)

## 🚀 Running PostgreSQL with Docker

### Option 1: Using Docker Compose (Recommended)

**Start the database:**
```bash
docker-compose up -d postgres
```

**Check status:**
```bash
docker ps | grep ecommerce
```

**View logs:**
```bash
docker logs ecommerce_postgres -f
```

**Stop the database:**
```bash
docker-compose down
```

**Reset database (delete all data):**
```bash
docker-compose down -v
docker-compose up -d postgres
```

### Option 2: Using Makefile Commands

```bash
# Start PostgreSQL
make db-start

# Stop PostgreSQL  
make db-stop

# Reset database (removes all data)
make db-reset

# View database logs
make db-logs

# Access PostgreSQL shell
make db-shell

# Check database status
make db-status
```

### Option 3: Using Helper Scripts

```bash
# Start database
./scripts/start-db.sh

# Stop database
./scripts/stop-db.sh

# Reset database
./scripts/reset-db.sh
```

## 🔧 Database Configuration

The database is configured to run on **port 5433** to avoid conflicts with other PostgreSQL installations.

### Connection Details

- **Host:** localhost
- **Port:** 5433 (mapped from container's 5432)
- **Database:** ecommerce_db
- **Username:** postgres
- **Password:** postgres
- **Connection String:** `postgresql://postgres:postgres@localhost:5433/ecommerce_db`

## 📝 Environment Variables

The `.env` file is already configured:

```env
DB_HOST=localhost
DB_PORT=5433
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=ecommerce_db
DB_SSLMODE=disable
```

## 🏃 Running the Application

Once PostgreSQL is running with Docker:

```bash
# Run the API
go run cmd/api/main.go

# Or using Makefile
make run

# Or build and run
make build
./bin/api
```

## 🧪 Testing the Setup

1. **Start PostgreSQL:**
   ```bash
   docker-compose up -d postgres
   ```

2. **Wait for it to be ready (3-5 seconds):**
   ```bash
   sleep 3
   ```

3. **Verify PostgreSQL is running:**
   ```bash
   docker exec ecommerce_postgres pg_isready -U postgres -d ecommerce_db
   ```
   
   Should output: `/var/run/postgresql:5432 - accepting connections`

4. **Run the API:**
   ```bash
   go run cmd/api/main.go
   ```

5. **Test the API:**
   ```bash
   curl http://localhost:8080/health
   ```
   
   Should output: `{"status":"ok"}`

## 🗄️ Database Operations

### Access PostgreSQL Shell

```bash
docker exec -it ecommerce_postgres psql -U postgres -d ecommerce_db
```

Common SQL commands:
```sql
-- List all tables
\dt

-- Describe a table
\d users

-- View all users
SELECT * FROM users;

-- Exit
\q
```

### Manual Database Creation (if needed)

```bash
# Access PostgreSQL
docker exec -it ecommerce_postgres psql -U postgres

# Create database
CREATE DATABASE ecommerce_db;

# Connect to database
\c ecommerce_db

# Exit
\q
```

### Backup Database

```bash
docker exec ecommerce_postgres pg_dump -U postgres ecommerce_db > backup.sql
```

### Restore Database

```bash
cat backup.sql | docker exec -i ecommerce_postgres psql -U postgres -d ecommerce_db
```

## 🛠️ Troubleshooting

### Port Already in Use

If port 5433 is already in use, you can change it:

1. Edit `docker-compose.yml`:
   ```yaml
   ports:
     - "5434:5432"  # Use port 5434 instead
   ```

2. Update `.env`:
   ```env
   DB_PORT=5434
   ```

3. Restart:
   ```bash
   docker-compose down
   docker-compose up -d postgres
   ```

### Database Connection Failed

**Check if container is running:**
```bash
docker ps | grep ecommerce_postgres
```

**Check logs for errors:**
```bash
docker logs ecommerce_postgres
```

**Verify environment variables:**
```bash
cat .env
```

### Reset Everything

If something goes wrong, reset everything:

```bash
# Stop and remove everything
docker-compose down -v

# Remove container if it still exists
docker rm -f ecommerce_postgres

# Remove volume if it still exists
docker volume rm crud-ecommerce_postgres_data

# Start fresh
docker-compose up -d postgres
```

### Can't Connect from Application

1. Verify PostgreSQL is running:
   ```bash
   docker ps | grep ecommerce
   ```

2. Test connection:
   ```bash
   docker exec ecommerce_postgres pg_isready -U postgres -d ecommerce_db
   ```

3. Check .env file matches docker-compose.yml

4. Restart the API

## 📊 Docker Compose File Structure

The `docker-compose.yml` defines:

- **Service:** PostgreSQL 15 Alpine
- **Container Name:** ecommerce_postgres
- **Ports:** 5433:5432 (host:container)
- **Volume:** postgres_data (persistent storage)
- **Health Check:** Ensures database is ready before use
- **Auto-restart:** Unless manually stopped

## 🔄 Development Workflow

### Daily Development

```bash
# Morning: Start database
make db-start

# Work on your code
go run cmd/api/main.go

# Evening: Stop database (optional)
make db-stop
```

### Testing Changes

```bash
# Reset database to clean state
make db-reset

# Run application
go run cmd/api/main.go

# Test your endpoints
curl http://localhost:8080/api/v1/...
```

## 🎯 Production Considerations

For production, you should:

1. **Use strong passwords:**
   ```yaml
   environment:
     POSTGRES_PASSWORD: ${POSTGRES_PASSWORD}  # From secrets
   ```

2. **Enable SSL:**
   ```env
   DB_SSLMODE=require
   ```

3. **Use managed database:**
   - AWS RDS
   - Google Cloud SQL
   - Azure Database for PostgreSQL
   - DigitalOcean Managed Databases

4. **Set resource limits:**
   ```yaml
   deploy:
     resources:
       limits:
         cpus: '1'
         memory: 1G
   ```

5. **Regular backups:**
   ```bash
   # Automated backup script
   docker exec ecommerce_postgres pg_dump -U postgres ecommerce_db | gzip > backup_$(date +%Y%m%d).sql.gz
   ```

## 📚 Additional Resources

- [Docker Documentation](https://docs.docker.com/)
- [PostgreSQL Docker Image](https://hub.docker.com/_/postgres)
- [Docker Compose Documentation](https://docs.docker.com/compose/)

## ✅ Checklist

- [x] Docker installed and running
- [x] PostgreSQL container started
- [x] Database created (ecommerce_db)
- [x] Application connects successfully
- [x] Migrations run successfully
- [x] API responds to health check

---

**Happy Coding! 🚀**
