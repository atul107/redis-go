package handlers

import (
	. "github.com/redis-go/internal/aof"
	. "github.com/redis-go/pkg/resp"
)

// HgetHandler handles the "HGET" command.
func HgetHandler(value Value, aof *Aof) Value {
	args := value.Array[1:]

	if len(args) != 2 {
		return Value{Typ: "error", Str: "ERR wrong number of arguments for 'hget' command"}
	}

	hash := args[0].Bulk
	key := args[1].Bulk

	HashStoreLock.RLock()
	res, ok := HashStore[hash][key]
	HashStoreLock.RUnlock()

	if !ok {
		return Value{Typ: "null"}
	}

	return Value{Typ: "bulk", Bulk: res}
}
