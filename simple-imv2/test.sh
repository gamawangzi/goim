#!/bin/bash

echo "=== 微服务架构测试 ==="
echo ""

echo "1. 用户 user123 加入房间 room001"
curl -X POST http://localhost:3111/room/join \
  -H "Content-Type: application/json" \
  -d '{"user_id":"user123","room_id":"room001"}'
echo -e "\n"

sleep 1

echo "2. 用户 user456 加入房间 room001"
curl -X POST http://localhost:3111/room/join \
  -H "Content-Type: application/json" \
  -d '{"user_id":"user456","room_id":"room001"}'
echo -e "\n"

sleep 1

echo "3. 向房间 room001 推送消息"
