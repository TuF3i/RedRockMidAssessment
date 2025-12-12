local stockKey  = KEYS[1] -- 课程余量的键
local usersKey  = KEYS[2] -- 该课程的学生表键
local userSelectedKey = KEYS[3] -- 学生自己的选课表键

local userID    = ARGV[1]
local courseID  = ARGV[2]

-- 减库存
local left = redis.call('DECR', stockKey)
if left < 0 then
-- 库存不足，回滚刚才的 DECR，返回 0
	redis.call('INCR', stockKey)
	return 0
end

-- 把用户加入集合
redis.call('SADD', usersKey, userID)
redis.call('SADD', userSelectedKey, courseID)

return 1