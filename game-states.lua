-- local source_key = '01-choosing-trump'
local source_key = '02-first-card-has-placed'
local destination_key = 'game:f2eed7f0-dc65-457d-aeac-54a35ffde75b'
local fields_values = redis.call('HGETALL', source_key)
if next(fields_values) then
    redis.call('HMSET', destination_key, unpack(fields_values))
end

