# 🚀 Local Development Setup Guide

## 📋 Step 1: Check What You Already Have

Run these commands to check what's already installed:

```bash
# Check Go
go version

# Check Docker
docker --version
docker-compose --version

# Check Redis (if installed)
redis-cli --version

# Check PostgreSQL (if installed)
psql --version

# Check Git
git --version
```

## 🔧 Step 2: Install Required Tools

### A. Install Go (Required)

#### Windows:
1. Download from: https://golang.org/dl/
2. Run installer and follow instructions
3. Verify: `go version`

#### macOS:
```bash
# Using Homebrew (recommended)
brew install go

# Or download from: https://golang.org/dl/
```

#### Linux (Ubuntu/Debian):
```bash
# Method 1: Using package manager
sudo apt update
sudo apt install golang-go

# Method 2: Latest version
wget https://go.dev/dl/go1.21.0.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.21.0.linux-amd64.tar.gz
export PATH=$PATH:/usr/local/go/bin
```

### B. Install Docker & Docker Compose (Recommended)

#### Windows:
1. Download Docker Desktop: https://www.docker.com/products/docker-desktop
2. Install and restart
3. Verify: `docker --version`

#### macOS:
```bash
# Using Homebrew
brew install --cask docker

# Or download Docker Desktop from: https://www.docker.com/products/docker-desktop
```

#### Linux (Ubuntu/Debian):
```bash
# Install Docker
sudo apt update
sudo apt install docker.io docker-compose
sudo systemctl start docker
sudo systemctl enable docker
sudo usermod -aG docker $USER

# Log out and log back in, then verify
docker --version
docker-compose --version
```

### C. Install Redis (Optional - for direct testing)

#### Windows:
```bash
# Using WSL2 or download from: https://github.com/microsoftarchive/redis/releases
```

#### macOS:
```bash
brew install redis
brew services start redis
```

#### Linux (Ubuntu/Debian):
```bash
sudo apt update
sudo apt install redis-server
sudo systemctl start redis-server
sudo systemctl enable redis-server
```

### D. Install PostgreSQL (Optional - for direct testing)

#### Windows:
1. Download from: https://www.postgresql.org/download/windows/
2. Run installer with default settings

#### macOS:
```bash
brew install postgresql
brew services start postgresql
```

#### Linux (Ubuntu/Debian):
```bash
sudo apt update
sudo apt install postgresql postgresql-contrib
sudo systemctl start postgresql
sudo systemctl enable postgresql
```

## 🎯 Step 3: Quick Setup Options

### Option 1: Docker Compose (Easiest - Recommended)
```bash
# Navigate to project
cd core-services/egov-url-shortening-go

# Start everything (includes PostgreSQL, Redis, and the app)
docker-compose up --build

# Access service at: http://localhost:8091
```

### Option 2: Local Go with Docker Dependencies
```bash
# Start only databases with Docker
docker run -d --name redis -p 6379:6379 redis:7-alpine
docker run -d --name postgres -p 5432:5432 -e POSTGRES_PASSWORD=postgres postgres:15-alpine

# Run Go service locally
go mod tidy
go run main.go
```

### Option 3: Full Local Setup
```bash
# Start Redis
redis-server

# Start PostgreSQL (in another terminal)
sudo -u postgres psql -c "CREATE DATABASE devdb;"

# Run Go service
export DATABASE_ENABLED=true
export DATABASE_HOST=localhost
export DATABASE_PASSWORD=postgres
go run main.go
```

### Option 4: Memory-Only Testing (No external dependencies)
```bash
# Just run with in-memory storage
export DATABASE_ENABLED=false
export REDIS_HOST=unavailable  # This will trigger memory fallback
go run main.go
```

## 🧪 Step 4: Verify Installation

After choosing an option above, test the service:

```bash
# Health check
curl http://localhost:8091/health

# Shorten a URL
curl -X POST http://localhost:8091/egov-url-shortening/shortener \
  -H "Content-Type: application/json" \
  -d '{"url": "https://www.google.com"}'

# Test redirect (use the returned shortened URL)
curl -L http://localhost:8091/egov-url-shortening/{returned-id}
```

## 🚨 Troubleshooting

### Common Issues:

1. **Port already in use**:
   ```bash
   # Find what's using port 8091
   lsof -i :8091
   # Kill the process or change port in config
   ```

2. **Permission denied (Docker)**:
   ```bash
   sudo usermod -aG docker $USER
   # Log out and log back in
   ```

3. **Go modules not found**:
   ```bash
   go mod tidy
   go mod download
   ```

4. **Database connection failed**:
   ```bash
   # Check if services are running
   docker ps
   # Or check local services
   systemctl status redis postgresql
   ```

## 📱 Step 5: Postman Testing

### Import Collection:
Create a new Postman collection with these requests:

1. **Health Check**
   - Method: GET
   - URL: `http://localhost:8091/health`

2. **Shorten URL**
   - Method: POST
   - URL: `http://localhost:8091/egov-url-shortening/shortener`
   - Headers: `Content-Type: application/json`
   - Body: `{"url": "https://www.example.com"}`

3. **Get URL Details**
   - Method: GET
   - URL: `http://localhost:8091/egov-url-shortening/details/{id}`

4. **Delete URL**
   - Method: DELETE
   - URL: `http://localhost:8091/egov-url-shortening/{id}`

5. **Service Stats**
   - Method: GET
   - URL: `http://localhost:8091/stats`

## 🎉 Next Steps

Once everything is working:
1. ✅ Test all endpoints with Postman
2. ✅ Run load tests if needed
3. ✅ Deploy to your target environment
4. ✅ Update Java code to match Go structure

Choose your preferred setup method and let me know if you need help with any specific step!