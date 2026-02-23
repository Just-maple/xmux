#!/bin/bash

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${YELLOW}Starting Gin Business Example API Test${NC}"
echo "========================================"

# Start server in background
echo -e "${YELLOW}Starting server...${NC}"
go run cmd/server/main.go &
SERVER_PID=$!

# Wait for server to start
sleep 3

echo -e "${YELLOW}Testing API endpoints...${NC}"

# Test 1: Register a new user
echo -e "\n${YELLOW}Test 1: Register User${NC}"
REGISTER_RESPONSE=$(curl -s -X POST http://localhost:8080/api/v1/users/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "email": "test@example.com",
    "password": "TestPass123!",
    "full_name": "Test User",
    "role": "user"
  }')

if echo "$REGISTER_RESPONSE" | grep -q "testuser"; then
  echo -e "${GREEN}✓ User registered successfully${NC}"
else
  echo -e "${RED}✗ User registration failed${NC}"
  echo "Response: $REGISTER_RESPONSE"
fi

# Test 2: Login
echo -e "\n${YELLOW}Test 2: Login${NC}"
LOGIN_RESPONSE=$(curl -s -X POST http://localhost:8080/api/v1/users/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "password": "TestPass123!"
  }')

if echo "$LOGIN_RESPONSE" | grep -q "token"; then
  echo -e "${GREEN}✓ Login successful${NC}"
  # Extract token for subsequent requests
  TOKEN=$(echo "$LOGIN_RESPONSE" | grep -o '"token":"[^"]*"' | cut -d'"' -f4)
else
  echo -e "${RED}✗ Login failed${NC}"
  echo "Response: $LOGIN_RESPONSE"
fi

# Test 3: Get user profile (if we have a token)
if [ -n "$TOKEN" ]; then
  echo -e "\n${YELLOW}Test 3: Get User Profile${NC}"
  PROFILE_RESPONSE=$(curl -s -X GET http://localhost:8080/api/v1/users/me \
    -H "Authorization: Bearer $TOKEN")
  
  if echo "$PROFILE_RESPONSE" | grep -q "testuser"; then
    echo -e "${GREEN}✓ Profile retrieved successfully${NC}"
  else
    echo -e "${RED}✗ Profile retrieval failed${NC}"
    echo "Response: $PROFILE_RESPONSE"
  fi
fi

# Test 4: List users (admin required)
echo -e "\n${YELLOW}Test 4: List Users${NC}"
LIST_RESPONSE=$(curl -s -X GET "http://localhost:8080/api/v1/users?limit=5&offset=0")

if echo "$LIST_RESPONSE" | grep -q "users"; then
  echo -e "${GREEN}✓ Users listed successfully${NC}"
else
  echo -e "${RED}✗ User listing failed${NC}"
  echo "Response: $LIST_RESPONSE"
fi

# Stop server
echo -e "\n${YELLOW}Stopping server...${NC}"
kill $SERVER_PID 2>/dev/null
wait $SERVER_PID 2>/dev/null

echo -e "\n${YELLOW}Test completed${NC}"