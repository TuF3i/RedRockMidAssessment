local stockKey   = KEYS[1]   -- 课程库存
local usersKey   = KEYS[2]   -- 已报名用户集合
local droppedKey = KEYS[3]   -- 已退课用户集合
local userSelectedKey = KEYS[4] -- 学生选课集合

local userID     = ARGV[1]   -- 要退课的用户
local courseID   = ARGV[2]   -- 课程ID

-- 1. 减库存（退课逻辑里也可 INCR，这里演示 DECR 场景）
local left = redis.call('INCR', stockKey)

-- 2. 从报名集合移除
if redis.call('EXISTS', usersKey) == 1 then
	redis.call('SREM', usersKey, userID)
end

-- 3. 加入退课集合
redis.call('SADD', droppedKey, userID)

-- 4. 从学生选课集合中移除课程
if redis.call('EXISTS', userSelectedKey) == 1 then
	redis.call('SREM', userSelectedKey, courseID)
end

return 1