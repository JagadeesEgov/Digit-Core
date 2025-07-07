# 🚀 Getting Started - URL Shortener Service

## 🎯 Instant Test (No Installation Required)

If you just want to test the service immediately without installing anything:

```bash
cd core-services/egov-url-shortening-go
./quick-test.sh
```

This will:
- ✅ Run the service with in-memory storage
- ✅ Test all API endpoints automatically
- ✅ Show you exactly how everything works
- ✅ No Redis or PostgreSQL needed!

## 🛠 Step-by-Step Setup

### 1. Check What You Have
```bash
./check-setup.sh
```

This script will tell you exactly what you need to install.

### 2. Install Go (Required)

**Ubuntu/Debian:**
```bash
sudo apt update && sudo apt install golang-go
```

**macOS:**
```bash
brew install go
```

**Windows:** Download from https://golang.org/dl/

### 3. Choose Your Setup Method

#### Option A: Memory-Only (Easiest)
```bash
export DATABASE_ENABLED=false
go mod tidy
go run main.go
```
- ✅ No external dependencies
- ❌ Data doesn't persist between restarts

#### Option B: Docker Compose (Recommended)
```bash
# Install Docker first, then:
docker-compose up --build
```
- ✅ Includes Redis + PostgreSQL + Service
- ✅ Production-like environment
- ✅ Data persists

#### Option C: Local Databases
```bash
# Install Redis and PostgreSQL locally, then:
export DATABASE_ENABLED=true
export DATABASE_HOST=localhost
go run main.go
```

### 4. Test the Service

Once running, test with curl:

```bash
# Health check
curl http://localhost:8091/health

# Shorten a URL
curl -X POST http://localhost:8091/egov-url-shortening/shortener \
  -H "Content-Type: application/json" \
  -d '{"url": "https://www.google.com"}'

# Use the returned short URL to test redirect
```

## 📱 API Endpoints

| Method | URL | Description |
|--------|-----|-------------|
| GET | `/health` | Health check |
| POST | `/egov-url-shortening/shortener` | Create short URL |
| GET | `/egov-url-shortening/{id}` | Redirect to original URL |
| GET | `/egov-url-shortening/details/{id}` | Get URL details |
| DELETE | `/egov-url-shortening/{id}` | Delete URL |
| GET | `/stats` | Service statistics |

## 🔧 Configuration

Create a `.env` file or set environment variables:

```bash
# Database settings
DATABASE_ENABLED=true
DATABASE_HOST=localhost
DATABASE_PORT=5432
DATABASE_NAME=devdb
DATABASE_USER=postgres
DATABASE_PASSWORD=postgres

# Redis settings (fallback storage)
REDIS_HOST=localhost
REDIS_PORT=6379

# Server settings
SERVER_PORT=8091
SERVER_CONTEXT_PATH=/egov-url-shortening

# HashID settings (must match Java version)
HASHIDS_SALT=egov-url-shortening-salt
HASHIDS_MIN_LENGTH=5
```

## 🐳 Production Deployment

### Docker
```bash
# Build image
docker build -t url-shortener .

# Run with external databases
docker run -p 8091:8091 \
  -e DATABASE_HOST=your-db-host \
  -e REDIS_HOST=your-redis-host \
  url-shortener
```

### Binary Deployment
```bash
# Build for your target platform
go build -o url-shortener main.go

# Deploy binary + config
./url-shortener
```

## 🚨 Troubleshooting

**Port already in use:**
```bash
lsof -i :8091  # Find what's using the port
```

**Build errors:**
```bash
go mod tidy
go mod download
```

**Database connection failed:**
```bash
# Use memory mode for testing
export DATABASE_ENABLED=false
export REDIS_HOST=unavailable
```

## 🎉 Next Steps

1. ✅ Test locally with `./quick-test.sh`
2. ✅ Set up with your preferred database
3. ✅ Use Postman to test all endpoints
4. ✅ Deploy to your environment
5. ✅ Update Java code to match structure

## 📞 Need Help?

- 📖 Detailed setup: `setup-guide.md`
- 🔍 Check your system: `./check-setup.sh`
- ⚡ Quick test: `./quick-test.sh`
- 📊 Project overview: `PROJECT_SUMMARY.md`