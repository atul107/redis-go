package handlers

import (
	. "github.com/redis-go/internal/aof"
	. "github.com/redis-go/pkg/resp"
)

func DelHandler(value Value, aof *Aof) Value {
	args := value.Array[1:]

	if len(args) < 1 {
		return Value{Typ: "error", Str: "ERR wrong number of arguments for 'del' command"}
	}

	keys := make([]string, len(args))
	for i, arg := range args {
		keys[i] = arg.Bulk
	}

	deletedKeys := []string{}

	KeyValueStoreLock.Lock()
	for _, key := range keys {
		if _, ok := KeyValueStore[key]; ok {
			delete(KeyValueStore, key)
			deletedKeys = append(deletedKeys, key)
		}
	}
	KeyValueStoreLock.Unlock()

	for _, key := range deletedKeys {
		aof.Write(Value{Typ: "array", Array: []Value{
			{Typ: "bulk", Bulk: "DEL"},
			{Typ: "bulk", Bulk: key},
		}})
	}
	return Value{Typ: "integer", Num: len(deletedKeys)}
}
