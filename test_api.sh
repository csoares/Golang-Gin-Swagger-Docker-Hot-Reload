#!/bin/bash

# API Test Script
# Tests all endpoints of the Go/Gin API

BASE_URL="http://localhost:8081"
ECHO_ENDPOINT="$BASE_URL/api/v1/echo"
LOGIN_ENDPOINT="$BASE_URL/api/v1/auth/login"
REGISTER_ENDPOINT="$BASE_URL/api/v1/auth/register"
REFRESH_ENDPOINT="$BASE_URL/api/v1/auth/refresh_token"
EVALUATION_ENDPOINT="$BASE_URL/api/v1/evaluation"

PASS_COUNT=0
FAIL_COUNT=0

# Main execution
echo "=========================================="
echo "       API Endpoint Test Suite"
echo "=========================================="
echo ""

# Test 1: Echo endpoint
echo -n "1. Testing GET /api/v1/echo ... "
RESPONSE=$(curl -s "$ECHO_ENDPOINT?name=test" 2>/dev/null)
if echo "$RESPONSE" | grep -q '"echo":"test"'; then
    echo "✓ PASS"
    ((PASS_COUNT++))
else
    echo "✗ FAIL - Response: $RESPONSE"
    ((FAIL_COUNT++))
fi

# Get auth token
echo -n "2. Testing POST /api/v1/auth/login ... "
LOGIN_RESPONSE=$(curl -s -X POST "$LOGIN_ENDPOINT" \
    -H "Content-Type: application/json" \
    -d '{"username":"admin","password":"admin123"}' 2>/dev/null)
TOKEN=$(echo "$LOGIN_RESPONSE" | grep -o '"token":"[^"]*"' | cut -d'"' -f4)
if [ -n "$TOKEN" ]; then
    echo "✓ PASS"
    ((PASS_COUNT++))
else
    echo "✗ FAIL - Response: $LOGIN_RESPONSE"
    ((FAIL_COUNT++))
fi

# Test register
echo -n "3. Testing POST /api/v1/auth/register ... "
RANDOM_USER="user_$(date +%s)"
REGISTER_RESPONSE=$(curl -s -X POST "$REGISTER_ENDPOINT" \
    -H "Content-Type: application/json" \
    -d "{\"username\":\"$RANDOM_USER\",\"password\":\"password123\"}" 2>/dev/null)
if echo "$REGISTER_RESPONSE" | grep -q '"status":201'; then
    echo "✓ PASS (created user: $RANDOM_USER)"
    ((PASS_COUNT++))
else
    echo "✗ FAIL - Response: $REGISTER_RESPONSE"
    ((FAIL_COUNT++))
fi

# Test create evaluation
echo -n "4. Testing POST /api/v1/evaluation/ ... "
CREATE_RESPONSE=$(curl -s -X POST "$EVALUATION_ENDPOINT/" \
    -H "Authorization: Bearer $TOKEN" \
    -H "Content-Type: application/json" \
    -d '{"Rating":5,"Note":"Test evaluation"}' 2>/dev/null)
NEW_ID=$(echo "$CREATE_RESPONSE" | grep -o '"resourceId":[0-9]*' | cut -d':' -f2)
if [ -n "$NEW_ID" ]; then
    echo "✓ PASS (created ID: $NEW_ID)"
    ((PASS_COUNT++))
else
    echo "✗ FAIL - Response: $CREATE_RESPONSE"
    ((FAIL_COUNT++))
fi

# Test get all evaluations
echo -n "5. Testing GET /api/v1/evaluation/ ... "
LIST_RESPONSE=$(curl -s "$EVALUATION_ENDPOINT/" \
    -H "Authorization: Bearer $TOKEN" 2>/dev/null)
if echo "$LIST_RESPONSE" | grep -q '"status":200'; then
    COUNT=$(echo "$LIST_RESPONSE" | grep -o '"id"' | wc -l | tr -d ' ')
    echo "✓ PASS ($COUNT evaluations found)"
    ((PASS_COUNT++))
else
    echo "✗ FAIL - Response: $LIST_RESPONSE"
    ((FAIL_COUNT++))
fi

# Test get by ID
echo -n "6. Testing GET /api/v1/evaluation/$NEW_ID ... "
GET_RESPONSE=$(curl -s "$EVALUATION_ENDPOINT/$NEW_ID" \
    -H "Authorization: Bearer $TOKEN" 2>/dev/null)
if echo "$GET_RESPONSE" | grep -q '"status":200'; then
    echo "✓ PASS"
    ((PASS_COUNT++))
else
    echo "✗ FAIL - Response: $GET_RESPONSE"
    ((FAIL_COUNT++))
fi

# Test update
echo -n "7. Testing PUT /api/v1/evaluation/$NEW_ID ... "
UPDATE_RESPONSE=$(curl -s -X PUT "$EVALUATION_ENDPOINT/$NEW_ID" \
    -H "Authorization: Bearer $TOKEN" \
    -H "Content-Type: application/json" \
    -d '{"Rating":4,"Note":"Updated evaluation"}' 2>/dev/null)
if echo "$UPDATE_RESPONSE" | grep -q '"status":200'; then
    echo "✓ PASS"
    ((PASS_COUNT++))
else
    echo "✗ FAIL - Response: $UPDATE_RESPONSE"
    ((FAIL_COUNT++))
fi

# Test raw query
echo -n "8. Testing GET /api/v1/evaluation/raw ... "
RAW_RESPONSE=$(curl -s "$EVALUATION_ENDPOINT/raw?min_rating=4" \
    -H "Authorization: Bearer $TOKEN" 2>/dev/null)
if echo "$RAW_RESPONSE" | grep -q '"status":200'; then
    echo "✓ PASS"
    ((PASS_COUNT++))
else
    echo "✗ FAIL - Response: $RAW_RESPONSE"
    ((FAIL_COUNT++))
fi

# Test batch update
echo -n "9. Testing PUT /api/v1/evaluation/batch ... "
BATCH_RESPONSE=$(curl -s -X PUT "$EVALUATION_ENDPOINT/batch" \
    -H "Authorization: Bearer $TOKEN" \
    -H "Content-Type: application/json" \
    -d '{"increment":1}' 2>/dev/null)
if echo "$BATCH_RESPONSE" | grep -q '"status":200'; then
    echo "✓ PASS"
    ((PASS_COUNT++))
else
    echo "✗ FAIL - Response: $BATCH_RESPONSE"
    ((FAIL_COUNT++))
fi

# Test create with audit
echo -n "10. Testing POST /api/v1/evaluation/audit ... "
AUDIT_RESPONSE=$(curl -s -X POST "$EVALUATION_ENDPOINT/audit" \
    -H "Authorization: Bearer $TOKEN" \
    -H "Content-Type: application/json" \
    -d '{"Rating":3,"Note":"Audit test"}' 2>/dev/null)
if echo "$AUDIT_RESPONSE" | grep -q '"status":201'; then
    echo "✓ PASS"
    ((PASS_COUNT++))
else
    echo "✗ FAIL - Response: $AUDIT_RESPONSE"
    ((FAIL_COUNT++))
fi

# Test delete
echo -n "11. Testing DELETE /api/v1/evaluation/$NEW_ID ... "
DELETE_RESPONSE=$(curl -s -X DELETE "$EVALUATION_ENDPOINT/$NEW_ID" \
    -H "Authorization: Bearer $TOKEN" 2>/dev/null)
if echo "$DELETE_RESPONSE" | grep -q '"status":200'; then
    echo "✓ PASS"
    ((PASS_COUNT++))
else
    echo "✗ FAIL - Response: $DELETE_RESPONSE"
    ((FAIL_COUNT++))
fi

# Test refresh token
echo -n "12. Testing PUT /api/v1/auth/refresh_token ... "
REFRESH_RESPONSE=$(curl -s -X PUT "$REFRESH_ENDPOINT" \
    -H "Authorization: Bearer $TOKEN" 2>/dev/null)
if echo "$REFRESH_RESPONSE" | grep -q '"status":200'; then
    echo "✓ PASS"
    ((PASS_COUNT++))
else
    echo "✗ FAIL - Response: $REFRESH_RESPONSE"
    ((FAIL_COUNT++))
fi

echo ""
echo "=========================================="
echo "              Test Summary"
echo "=========================================="
echo "Passed: $PASS_COUNT"
echo "Failed: $FAIL_COUNT"
echo "=========================================="

if [ $FAIL_COUNT -eq 0 ]; then
    echo "All tests passed!"
    exit 0
else
    echo "Some tests failed."
    exit 1
fi
