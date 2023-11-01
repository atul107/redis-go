package handlers

import (
	"path/filepath"

	. "github.com/redis-go/internal/aof"
	. "github.com/redis-go/pkg/resp"
)

func KeysHandler(value Value, aof *Aof) Value {
	args := value.Array[1:]

	if len(args) != 1 {
		return Value{Typ: "error", Str: "ERR wrong number of arguments for 'keys' command"}
	}

	pattern := args[0].Bulk
	matchingKeys := []string{}

	KeyValueStoreLock.RLock()
	for key := range KeyValueStore {
		if match, _ := filepath.Match(pattern, key); match {
			matchingKeys = append(matchingKeys, key)
		}
	}
	KeyValueStoreLock.RUnlock()

	values := make([]Value, len(matchingKeys))
	for i, key := range matchingKeys {
		values[i] = Value{Typ: "bulk", Bulk: key}
	}

	return Value{Typ: "array", Array: values}
}
