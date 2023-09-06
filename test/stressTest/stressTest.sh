#!/bin/bash
echo "压力测试1：获取视频流(未登录状态)"
wrk -t12 -c400 -d30s "http://10.201.83.51:8889/douyin/feed/"

echo "压力测试2：获取视频流(登录状态)"
wrk -t12 -c400 -d30s "http://10.201.83.51:8889/douyin/feed/?token=toni123456"

# echo "压力测试3：注册"
# wrk -t12 -c400 -d30s -H "Content-Type: application/x-www-form-urlencoded" -s register.lua "http://10.201.83.51:8889/douyin/user/login/"

echo "压力测试4：登录"
wrk -t12 -c400 -d30s -H "Content-Type: application/x-www-form-urlencoded" -s login.lua "http://10.201.83.51:8889/douyin/user/register/"

echo "压力测试5：获取用户信息(未登录状态)"
wrk -t12 -c400 -d30s "http://10.201.83.51:8889/douyin/user/?user_id=1"

echo "压力测试6：获取用户信息(登录状态)"
wrk -t12 -c400 -d30s "http://10.201.83.51:8889/douyin/user/?user_id=1&token=ghostfather1234567"

echo "压力测试7：发布列表"
wrk -t12 -c400 -d30s "http://10.201.83.51:8889/douyin/publish/list/?user_id=1&token=toni123456"

echo "压力测试8：喜欢列表"
wrk -t12 -c400 -d30s "http://10.201.83.51:8889/douyin/favorite/list/?user_id=1&token=toni123456"

echo "压力测试9：评论列表"
wrk -t12 -c400 -d30s "http://10.201.83.51:8889/douyin/comment/list/?video_id=7&token=ghostfather1234567"

echo "压力测试10：关注列表"
wrk -t12 -c400 -d30s "http://10.201.83.51:8889/douyin/relation/follow/list/?user_id=1&token=toni123456"

echo "压力测试11：粉丝列表"
wrk -t12 -c400 -d30s "http://10.201.83.51:8889/douyin/relation/follower/list/?user_id=1&token=toni123456"

echo "压力测试12：好友列表"
wrk -t12 -c400 -d30s "http://10.201.83.51:8889/douyin/relation/friend/list/?user_id=1&token=toni123456"

echo "压力测试13：聊天记录"
wrk -t12 -c400 -d30s "http://10.201.83.51:8889/douyin/message/chat/?to_user_id=2&token=toni123456"
