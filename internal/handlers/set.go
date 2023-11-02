package handlers

import (
	"strconv"
	"time"

	. "github.com/redis-go/internal/aof"
	. "github.com/redis-go/pkg/resp"
)

func SetHandler(val Value, aof *Aof) Value {
	args := val.Array[1:]

	if len(args) < 2 || len(args) > 7 {
		return Value{Typ: "error", Str: "ERR wrong number of arguments for 'set' command"}
	}

	key := args[0].Bulk
	value := args[1].Bulk
	getValue := false

	for i := 2; i < len(args); i++ {
		switch args[i].Bulk {
		case "NX":
			if _, exists := KeyValueStore[key]; exists {
				return Value{Typ: "null"}
			}
		case "XX":
			if _, exists := KeyValueStore[key]; !exists {
				return Value{Typ: "null"}
			}
		case "EX":
			if i+1 < len(args) {
				seconds, err := strconv.Atoi(args[i+1].Bulk)
				if err != nil {
					return Value{Typ: "error", Str: "ERR invalid argument for 'EX'"}
				}
				ExpireKey(key, seconds)
				i++
			} else {
				return Value{Typ: "error", Str: "ERR missing argument for 'EX'"}
			}
		case "PX":
			if i+1 < len(args) {
				milliseconds, err := strconv.Atoi(args[i+1].Bulk)
				if err != nil {
					return Value{Typ: "error", Str: "ERR invalid argument for 'PX'"}
				}
				ExpireKey(key, milliseconds/1000)
				i++
			} else {
				return Value{Typ: "error", Str: "ERR missing argument for 'PX'"}
			}
		case "EXAT":
			if i+1 < len(args) {
				unixTimeSeconds, err := strconv.Atoi(args[i+1].Bulk)
				if err != nil {
					return Value{Typ: "error", Str: "ERR invalid argument for 'EXAT'"}
				}
				SetAbsoluteExpiration(key, unixTimeSeconds)
				i++
			} else {
				return Value{Typ: "error", Str: "ERR missing argument for 'EXAT'"}
			}
		case "PXAT":
			if i+1 < len(args) {
				unixTimeMilliseconds, err := strconv.Atoi(args[i+1].Bulk)
				if err != nil {
					return Value{Typ: "error", Str: "ERR invalid argument for 'PXAT'"}
				}
				SetAbsoluteExpiration(key, unixTimeMilliseconds/1000)
				i++
			} else {
				return Value{Typ: "error", Str: "ERR missing argument for 'PXAT'"}
			}
		case "KEEPTTL":
			if _, exists := ExpiryStore[key]; exists {
				i++
			} else {
				return Value{Typ: "error", Str: "ERR key has no associated expiration"}
			}
		case "GET":
			getValue = true
		default:
			return Value{Typ: "error", Str: "ERR syntax error"}
		}
	}

	previousValue := Value{Typ: "null"}
	if getValue {
		if val, exists := KeyValueStore[key]; exists {
			previousValue = Value{Typ: "bulk", Bulk: val}
			return previousValue
		}
	}

	KeyValueStoreLock.Lock()
	KeyValueStore[key] = value
	aof.Write(val)
	KeyValueStoreLock.Unlock()

	return Value{Typ: "string", Str: "OK"}
}

func ExpireKey(key string, seconds int) {
	ExpirySToreLock.Lock()
	ExpiryStore[key] = time.Now().Add(time.Duration(seconds) * time.Second)
	ExpirySToreLock.Unlock()
}

func SetAbsoluteExpiration(key string, unixTimeSeconds int) {
	ExpirySToreLock.Lock()
	ExpiryStore[key] = time.Unix(int64(unixTimeSeconds), 0)
	ExpirySToreLock.Unlock()
}
