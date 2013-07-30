--[[
Locking algorithm:
	if the key exists
		if caller is owner of lock
			update ttl
			return true
		else return false
	if the key does not exist
		set owner
		set ttl
		return true
--]]

local call_owner = ARGV[1]
local ttl = ARGV[2]
if redis.call('EXISTS', KEYS[1]) == 1 then
	local lock_owner = redis.call('GET', KEYS[1])
	if lock_owner == call_owner then
		redis.call('EXPIRE', KEYS[1], ttl)
		return 1
	end
	return 0
else
	redis.call('SET', KEYS[1], call_owner)
	redis.call('EXPIRE', KEYS[1], ttl)
	return 1
end
