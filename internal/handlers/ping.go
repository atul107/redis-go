package handlers

import (
	. "github.com/redis-go/internal/aof"
	. "github.com/redis-go/pkg/resp"
)

// PingHandler handles the "PING" command.
func PingHandler(value Value, aof *Aof) Value {
	args := value.Array[1:]

	if len(args) == 0 {
		return Value{Typ: "string", Str: "PONG"}
	}

	return Value{Typ: "string", Str: args[0].Bulk}
}
