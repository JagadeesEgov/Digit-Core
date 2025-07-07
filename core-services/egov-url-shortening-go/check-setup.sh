#!/bin/bash

# Check Setup Script for URL Shortener Service
echo "🔍 Checking your system setup..."
echo "======================================="

# Check OS
OS=$(uname -s)
echo "📱 Operating System: $OS"

# Function to check if command exists
check_command() {
    if command -v "$1" &> /dev/null; then
        echo "✅ $1 is installed: $($1 --version 2>&1 | head -n1)"
        return 0
    else
        echo "❌ $1 is NOT installed"
        return 1
    fi
}

# Function to suggest installation
suggest_install() {
    echo "📦 Installation suggestions for $OS:"
    case $OS in
        "Linux")
            echo "   sudo apt update && sudo apt install $1"
            ;;
        "Darwin")
            echo "   brew install $2"
            ;;
        "MINGW"*|"CYGWIN"*|"MSYS"*)
            echo "   Download from: $3"
            ;;
        *)
            echo "   Please check the setup guide for your OS"
            ;;
    esac
}

# Check Go
echo -e "\n🔧 Checking Go..."
if ! check_command go; then
    suggest_install "golang-go" "go" "https://golang.org/dl/"
fi

# Check Docker
echo -e "\n🐳 Checking Docker..."
if ! check_command docker; then
    suggest_install "docker.io" "docker" "https://www.docker.com/products/docker-desktop"
fi

# Check Docker Compose
echo -e "\n🔄 Checking Docker Compose..."
if ! check_command docker-compose; then
    suggest_install "docker-compose" "docker-compose" "Included with Docker Desktop"
fi

# Check Redis (optional)
echo -e "\n🔴 Checking Redis (optional)..."
if ! check_command redis-cli; then
    echo "⚠️  Redis not found (optional - can use Docker)"
    suggest_install "redis-server" "redis" "https://github.com/microsoftarchive/redis/releases"
fi

# Check PostgreSQL (optional)
echo -e "\n🐘 Checking PostgreSQL (optional)..."
if ! check_command psql; then
    echo "⚠️  PostgreSQL not found (optional - can use Docker)"
    suggest_install "postgresql postgresql-contrib" "postgresql" "https://www.postgresql.org/download/"
fi

# Check curl (for testing)
echo -e "\n🌐 Checking curl (for testing)..."
check_command curl

echo -e "\n🎯 Setup Recommendations:"
echo "================================"

# Count installed tools
go_installed=$(command -v go &> /dev/null && echo 1 || echo 0)
docker_installed=$(command -v docker &> /dev/null && echo 1 || echo 0)

if [ "$go_installed" -eq 1 ] && [ "$docker_installed" -eq 1 ]; then
    echo "🚀 You have Go and Docker! You can use Option 1 or 2:"
    echo "   Option 1: docker-compose up --build"
    echo "   Option 2: Run databases in Docker, Go locally"
elif [ "$go_installed" -eq 1 ]; then
    echo "🟡 You have Go but no Docker. You can:"
    echo "   - Install local Redis/PostgreSQL"
    echo "   - Or use memory-only mode (no persistence)"
    echo "   Memory-only command: export DATABASE_ENABLED=false && go run main.go"
elif [ "$docker_installed" -eq 1 ]; then
    echo "🔵 You have Docker but no Go. Install Go first:"
    echo "   Then you can use Docker Compose option"
else
    echo "🟠 Install Go first for the easiest setup:"
    echo "   Go + memory-only mode needs no other dependencies"
fi

echo -e "\n💡 Quick Start (No dependencies needed):"
echo "cd core-services/egov-url-shortening-go"
echo "export DATABASE_ENABLED=false"
echo "export REDIS_HOST=unavailable"
echo "go mod tidy"
echo "go run main.go"

echo -e "\n📖 For detailed setup instructions, see: setup-guide.md"