--[[
Unlocking algorithm:
	if the key exists
		if owner of lock
			delete key
			return true
		else return false
	else return false
--]]
local call_owner = ARGV[1]
if redis.call('EXISTS', KEYS[1]) == 1 then
	local lock_owner = redis.call('GET', KEYS[1])
	if lock_owner == call_owner then
		redis.call('DEL', KEYS[1])
		return 1
	end
	return 0
else
	return 0
end
