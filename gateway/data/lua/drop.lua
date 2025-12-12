local stockKey   = KEYS[1]   -- 课程库存
local usersKey   = KEYS[2]   -- 已报名用户集合
local droppedKey = KEYS[3]   -- 已退课用户集合
local userSelectedKey = KEYS[4] -- 学生选课集合

local userID     = ARGV[1]   -- 要退课的用户
local courseID   = ARGV[2]   -- 课程ID

-- 1. 减库存（退课逻辑里也可 INCR，这里演示 DECR 场景）
local left = redis.call('DECR', stockKey)
if left < 0 then
-- 库存不允许为负，回滚
	redis.call('INCR', stockKey)
	return 0
end

-- 2. 从报名集合移除
redis.call('SREM', usersKey, userID)

-- 3. 加入退课集合
redis.call('SADD', droppedKey, userID)

-- 4. 从学生选课集合中移除课程
redis.call('SREM', userSelectedKey, courseID)

return 1