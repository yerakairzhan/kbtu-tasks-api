#!/bin/bash

echo "Task API Demo Script"
echo "===================="
echo ""

API_KEY="secret12345"
BASE_URL="http://localhost:8080/v1"

echo "1. Creating first task..."
curl -s -X POST -H "X-API-KEY: $API_KEY" -H "Content-Type: application/json" \
  -d '{"title":"Write unit tests"}' \
  $BASE_URL/tasks | jq .
echo ""

echo "2. Creating second task..."
curl -s -X POST -H "X-API-KEY: $API_KEY" -H "Content-Type: application/json" \
  -d '{"title":"Deploy service"}' \
  $BASE_URL/tasks | jq .
echo ""

echo "3. Getting all tasks..."
curl -s -H "X-API-KEY: $API_KEY" $BASE_URL/tasks | jq .
echo ""

echo "4. Getting task by ID (id=1)..."
curl -s -H "X-API-KEY: $API_KEY" $BASE_URL/tasks?id=1 | jq .
echo ""

echo "5. Updating task (mark as done)..."
curl -s -X PATCH -H "X-API-KEY: $API_KEY" -H "Content-Type: application/json" \
  -d '{"done":true}' \
  $BASE_URL/tasks?id=1 | jq .
echo ""

echo "6. Filtering tasks (done=true)..."
curl -s -H "X-API-KEY: $API_KEY" $BASE_URL/tasks?done=true | jq .
echo ""

echo "7. Filtering tasks (done=false)..."
curl -s -H "X-API-KEY: $API_KEY" $BASE_URL/tasks?done=false | jq .
echo ""

echo "8. Fetching external tasks..."
curl -s -H "X-API-KEY: $API_KEY" $BASE_URL/external-tasks | jq '. | length'
echo ""

echo "9. Deleting task (id=2)..."
curl -s -X DELETE -H "X-API-KEY: $API_KEY" $BASE_URL/tasks?id=2 | jq .
echo ""

echo "10. Getting all tasks after deletion..."
curl -s -H "X-API-KEY: $API_KEY" $BASE_URL/tasks | jq .
echo ""

echo "11. Testing invalid API key..."
curl -s -H "X-API-KEY: wrong" $BASE_URL/tasks | jq .
echo ""

echo "12. Testing invalid ID..."
curl -s -H "X-API-KEY: $API_KEY" $BASE_URL/tasks?id=999 | jq .
echo ""

echo "Demo completed!"