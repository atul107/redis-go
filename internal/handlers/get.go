package handlers

import (
	. "github.com/redis-go/internal/aof"
	. "github.com/redis-go/pkg/resp"
)

// GetHandler handles the "GET" command.
func GetHandler(value Value, aof *Aof) Value {
	args := value.Array[1:]

	if len(args) != 1 {
		return Value{Typ: "error", Str: "ERR wrong number of arguments for 'get' command"}
	}

	key := args[0].Bulk

	KeyValueStoreLock.RLock()
	res, ok := KeyValueStore[key]
	KeyValueStoreLock.RUnlock()

	if !ok {
		return Value{Typ: "null"}
	}

	return Value{Typ: "bulk", Bulk: res}
}
