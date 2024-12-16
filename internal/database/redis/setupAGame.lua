-- Input data
local gameID = KEYS[1] -- Game ID

-- Store points
redis.call('HSET', 'game:' .. gameID, 'points', cjson.encode({
    total = {a = 0, b = 0},
    current_round = {a = 0, b = 0}
}))

-- Store king_cards as an array of JSON objects
redis.call('HSET', 'game:' .. gameID, 'king_cards', cjson.encode({
    {player = "p1", card = "2S"},
    {player = "p2", card = "3H"}
}))

-- Store center_cards with the new format
redis.call('HSET', 'game:' .. gameID, 'center_cards', cjson.encode({
    {player = "p3", card = "10H"},
    {player = "p4", card = "3C"}
}))

-- Store groups
redis.call('HSET', 'game:' .. gameID, 'groups', cjson.encode({
    a = {"p1", "p2"},
    b = {"p3", "p4"}
}))

-- Store the current turn
redis.call('HSET', 'game:' .. gameID, 'current_turn', "p1")

-- Store cards
redis.call('HSET', 'game:' .. gameID, 'cards', cjson.encode({
    a = {"3h", "4s"},
    b = {"4c", "6h"}
}))

-- Store judge
redis.call('HSET', 'game:' .. gameID, 'judge', "p3")

-- Store trump
redis.call('HSET', 'game:' .. gameID, 'trump', "hearts")

return "Game updated successfully"
