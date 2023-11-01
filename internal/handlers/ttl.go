package handlers

import (
	"time"

	. "github.com/redis-go/internal/aof"
	. "github.com/redis-go/pkg/resp"
)

func TTLHandler(value Value, aof *Aof) Value {
	args := value.Array[1:]

	if len(args) != 1 {
		return Value{Typ: "error", Str: "ERR wrong Number of arguments for 'ttl' command"}
	}

	key := args[0].Bulk

	ExpiresMu.RLock()
	expirationTime, ok := Expires[key]
	ExpiresMu.RUnlock()

	if !ok {
		return Value{Typ: "integer", Num: -2} // Key does not exist.
	}

	remainingTime := int(expirationTime.Sub(time.Now()).Seconds())

	if remainingTime < 0 {
		return Value{Typ: "integer", Num: -1} // Key has no associated TTL.
	}

	return Value{Typ: "integer", Num: remainingTime}
}
