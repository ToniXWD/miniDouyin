-- post.lua
-- 设置请求方法为POST
wrk.method = "POST"
-- 设置请求体为JSON格式的数据
wrk.body = '{"username": "toni","content": "password, 123456"}'
-- 设置请求头为JSON格式
wrk.headers["Content-Type"] = "application/json"