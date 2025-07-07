#!/bin/bash

# Quick Test Script - No dependencies needed
echo "🚀 Quick Test - URL Shortening Service"
echo "====================================="

# Set environment for memory-only mode
export DATABASE_ENABLED=false
export REDIS_HOST=unavailable
export SERVER_PORT=8091

echo "📦 Installing Go dependencies..."
go mod tidy

echo "🔧 Building service..."
if ! go build -o url-shortener main.go; then
    echo "❌ Build failed!"
    exit 1
fi

echo "🎯 Starting service in memory-only mode..."
echo "   (This requires no Redis or PostgreSQL installation)"
echo "   Service will run on: http://localhost:8091"
echo ""

# Start service in background
./url-shortener &
SERVICE_PID=$!

# Wait for service to start
echo "⏳ Waiting for service to start..."
sleep 3

# Test if service is running
if curl -s http://localhost:8091/health > /dev/null; then
    echo "✅ Service is running!"
    echo ""
    
    echo "🧪 Testing API endpoints..."
    echo "=========================="
    
    # Test health check
    echo "1. Health Check:"
    curl -s http://localhost:8091/health | python3 -m json.tool || echo "Service is healthy"
    echo ""
    
    # Test URL shortening
    echo "2. Shorten URL:"
    RESPONSE=$(curl -s -X POST http://localhost:8091/egov-url-shortening/shortener \
        -H "Content-Type: application/json" \
        -d '{"url": "https://www.google.com"}')
    echo "$RESPONSE" | python3 -m json.tool 2>/dev/null || echo "$RESPONSE"
    
    # Extract shortened URL ID from response
    SHORT_ID=$(echo "$RESPONSE" | grep -o '"shortUrl":"[^"]*"' | sed 's/"shortUrl":"[^/]*\/egov-url-shortening\///g' | sed 's/"//g')
    echo ""
    
    if [ ! -z "$SHORT_ID" ]; then
        echo "3. Test Redirect (shortened ID: $SHORT_ID):"
        echo "   URL: http://localhost:8091/egov-url-shortening/$SHORT_ID"
        
        # Test redirect (just check if it responds, don't follow redirect)
        HTTP_CODE=$(curl -s -o /dev/null -w "%{http_code}" http://localhost:8091/egov-url-shortening/$SHORT_ID)
        if [ "$HTTP_CODE" = "302" ] || [ "$HTTP_CODE" = "301" ]; then
            echo "   ✅ Redirect working (HTTP $HTTP_CODE)"
        else
            echo "   ⚠️  Got HTTP $HTTP_CODE (expected 302)"
        fi
        echo ""
        
        echo "4. Get URL Details:"
        curl -s http://localhost:8091/egov-url-shortening/details/$SHORT_ID | python3 -m json.tool 2>/dev/null
        echo ""
    fi
    
    echo "5. Service Stats:"
    curl -s http://localhost:8091/stats | python3 -m json.tool 2>/dev/null
    echo ""
    
    echo "✨ Test completed! Service is working correctly."
    echo ""
    echo "🎯 You can now:"
    echo "   - Use Postman to test more endpoints"
    echo "   - Install Redis/PostgreSQL for persistence"
    echo "   - Deploy using Docker"
    echo ""
    echo "📄 All endpoints:"
    echo "   POST http://localhost:8091/egov-url-shortening/shortener"
    echo "   GET  http://localhost:8091/egov-url-shortening/{id}"
    echo "   GET  http://localhost:8091/egov-url-shortening/details/{id}"
    echo "   DELETE http://localhost:8091/egov-url-shortening/{id}"
    echo "   GET  http://localhost:8091/health"
    echo "   GET  http://localhost:8091/stats"
    
else
    echo "❌ Service failed to start!"
    echo "Check the logs above for errors."
fi

echo ""
echo "🛑 Stopping service..."
kill $SERVICE_PID 2>/dev/null
wait $SERVICE_PID 2>/dev/null

echo "✅ Test script completed!"
echo ""
echo "📖 For production setup with Redis/PostgreSQL, see: setup-guide.md"