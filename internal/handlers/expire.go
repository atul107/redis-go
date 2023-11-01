package handlers

import (
	"strconv"
	"sync"
	"time"

	. "github.com/redis-go/internal/aof"
	. "github.com/redis-go/pkg/resp"
)

var Expires = map[string]time.Time{}
var ExpiresMu = sync.RWMutex{}

func ExpireHandler(value Value, aof *Aof) Value {
	args := value.Array[1:]

	if len(args) != 2 {
		return Value{Typ: "error", Str: "ERR wrong number of arguments for 'expire' command"}
	}

	key := args[0].Bulk
	seconds, err := strconv.Atoi(args[1].Bulk)
	if err != nil {
		return Value{Typ: "error", Str: "ERR value is not an integer or out of range"}
	}

	ExpiresMu.Lock()
	Expires[key] = time.Now().Add(time.Duration(seconds) * time.Second)
	ExpiresMu.Unlock()

	return Value{Typ: "integer", Num: 1}
}
