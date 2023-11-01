package handlers

import (
	. "github.com/redis-go/internal/aof"
	. "github.com/redis-go/pkg/resp"
)

// SetHandler handles the "SET" command.
func SetHandler(value Value, aof *Aof) Value {
	args := value.Array[1:]

	if len(args) != 2 {
		return Value{Typ: "error", Str: "ERR wrong number of arguments for 'set' command"}
	}

	key := args[0].Bulk
	res := args[1].Bulk

	KeyValueStoreLock.Lock()
	KeyValueStore[key] = res
	aof.Write(value)
	KeyValueStoreLock.Unlock()

	return Value{Typ: "string", Str: "OK"}
}
