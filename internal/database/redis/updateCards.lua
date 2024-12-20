redis.call('HSET', 'game:' .. KEYS[1], 'cards', cjson.encode({
    [0] = cjson.decode(ARGV[1]),
    [1] = cjson.decode(ARGV[2]),
    [2] = cjson.decode(ARGV[3]),
    [3] = cjson.decode(ARGV[4])
}))