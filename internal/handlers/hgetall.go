package handlers

import (
	. "github.com/redis-go/internal/aof"
	. "github.com/redis-go/pkg/resp"
)

// HgetallHandler handles the "HGETALL" command.
func HgetallHandler(value Value, aof *Aof) Value {
	args := value.Array[1:]

	if len(args) != 1 {
		return Value{Typ: "error", Str: "ERR wrong number of arguments for 'hgetall' command"}
	}

	hash := args[0].Bulk

	HashStoreLock.RLock()
	res, ok := HashStore[hash]
	HashStoreLock.RUnlock()

	if !ok {
		return Value{Typ: "null"}
	}

	values := []Value{}
	for k, v := range res {
		values = append(values, Value{Typ: "bulk", Bulk: k})
		values = append(values, Value{Typ: "bulk", Bulk: v})
	}

	return Value{Typ: "array", Array: values}
}
