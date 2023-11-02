package handlers

import (
	"strconv"
	"time"

	. "github.com/redis-go/internal/aof"
	. "github.com/redis-go/pkg/resp"
)

func ExpireHandler(value Value, aof *Aof) Value {
	args := value.Array[1:]

	if len(args) < 2 || len(args) > 3 {
		return Value{Typ: "error", Str: "ERR wrong number of arguments for 'expire' command"}
	}

	key := args[0].Bulk
	seconds, err := strconv.Atoi(args[1].Bulk)
	if err != nil {
		return Value{Typ: "error", Str: "ERR value is not an integer or out of range"}
	}

	ExpirySToreLock.Lock()

	if len(args) == 3 {
		switch args[2].Bulk {
		case "NX":
			if _, exists := ExpiryStore[key]; !exists {
				ExpiryStore[key] = time.Now().Add(time.Duration(seconds) * time.Second)
			} else {
				ExpirySToreLock.Unlock()
				return Value{Typ: "integer", Num: 0}
			}
		case "XX":
			if _, exists := ExpiryStore[key]; exists {
				ExpiryStore[key] = time.Now().Add(time.Duration(seconds) * time.Second)
			} else {
				ExpirySToreLock.Unlock()
				return Value{Typ: "integer", Num: 0}
			}
		case "GT":
			if existingExpire, exists := ExpiryStore[key]; !exists || time.Now().Add(time.Duration(seconds)*time.Second).After(existingExpire) {
				ExpiryStore[key] = time.Now().Add(time.Duration(seconds) * time.Second)
			} else {
				ExpirySToreLock.Unlock()
				return Value{Typ: "integer", Num: 0}
			}
		case "LT":
			if existingExpire, exists := ExpiryStore[key]; !exists || time.Now().Add(time.Duration(seconds)*time.Second).Before(existingExpire) {
				ExpiryStore[key] = time.Now().Add(time.Duration(seconds) * time.Second)
			} else {
				ExpirySToreLock.Unlock()
				return Value{Typ: "integer", Num: 0}
			}
		default:
			ExpirySToreLock.Unlock()
			return Value{Typ: "error", Str: "ERR syntax error"}
		}
	} else {
		ExpiryStore[key] = time.Now().Add(time.Duration(seconds) * time.Second)
	}

	ExpirySToreLock.Unlock()

	return Value{Typ: "integer", Num: 1}
}
