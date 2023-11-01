package handlers

import (
	. "github.com/redis-go/internal/aof"
	. "github.com/redis-go/pkg/resp"
)

// HsetHandler handles the "HSET" command.
func HsetHandler(value Value, aof *Aof) Value {
	args := value.Array[1:]

	if len(args) != 3 {
		return Value{Typ: "error", Str: "ERR wrong number of arguments for 'hset' command"}
	}

	hash := args[0].Bulk
	key := args[1].Bulk
	res := args[2].Bulk

	HashStoreLock.Lock()
	if _, ok := HashStore[hash]; !ok {
		HashStore[hash] = make(map[string]string)
	}
	HashStore[hash][key] = res
	aof.Write(value)
	HashStoreLock.Unlock()

	return Value{Typ: "string", Str: "OK"}
}
