#!/bin/bash
echo "Limit + offset + order_by:"
curl -H "X-API-KEY: secret12345" "http://localhost:8080/v1/users?page=1&page_size=5&order_by=name"
echo "===================="
sleep 2

echo "Filtering by ID, Name, Email (repeat as needed for each field):"
curl -H "X-API-KEY: secret12345" "http://localhost:8080/v1/users?filter_id=user-10"
curl -H "X-API-KEY: secret12345" "http://localhost:8080/v1/users?name=Alice"
curl -H "X-API-KEY: secret12345" "http://localhost:8080/v1/users?email=alice1@example.com"
echo "===================="
sleep 2

echo "Filtering by three random fields plus pagination/order_by"
curl -H "X-API-KEY: secret12345" "http://localhost:8080/v1/users?page=2&page_size=5&order_by=birth_date&name=Grace&gender=female&birth_date=1994-07-12"
echo "===================="
sleep 2

echo "GetCommonFriends flow:"
curl -H "X-API-KEY: secret12345" "http://localhost:8080/v1/users/common-friends?user1_id=user-01&user2_id=user-02"
echo "===================="
sleep 2

echo "Database verification"
docker compose exec db psql -U postgres -d go_kbtu -c "SELECT COUNT(*) FROM users;"
docker compose exec db psql -U postgres -d go_kbtu -c "SELECT COUNT(*) FROM user_friends;"
echo "===================="
