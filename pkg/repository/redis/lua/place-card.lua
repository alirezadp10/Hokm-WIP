redis.call('HSET', "game:" .. KEYS[1], 'center_cards', ARGV[1])
redis.call('HSET', "game:" .. KEYS[1], 'lead_suit', ARGV[2])
redis.call('HSET', "game:" .. KEYS[1], 'who_has_won_the_cards', ARGV[3])
redis.call('HSET', "game:" .. KEYS[1], 'points', ARGV[4])
redis.call('HSET', "game:" .. KEYS[1], 'turn', ARGV[5])
redis.call('HSET', "game:" .. KEYS[1], 'king', ARGV[6])
redis.call('HSET', "game:" .. KEYS[1], 'was_the_king_changed', ARGV[7])
redis.call('HSET', "game:" .. KEYS[1], 'last_move_timestamp', ARGV[14])
redis.call('HSET', "game:" .. KEYS[1], 'trump', ARGV[15])
redis.call('HSET', "game:" .. KEYS[1], 'cards', cjson.encode({
    [0] = cjson.decode(ARGV[8]),
    [1] = cjson.decode(ARGV[9]),
    [2] = cjson.decode(ARGV[10]),
    [3] = cjson.decode(ARGV[11])
}))
redis.call('PUBLISH', "placing_card", KEYS[1] .. "|" .. ARGV[12] .. "|" .. ARGV[13])
return ""