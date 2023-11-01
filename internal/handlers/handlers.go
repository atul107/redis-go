package handlers

import (
	"sync"

	. "github.com/redis-go/internal/aof"
	. "github.com/redis-go/pkg/resp"
)

var Handlers = map[string]func(Value, *Aof) Value{
	"PING":    PingHandler,
	"SET":     SetHandler,
	"GET":     GetHandler,
	"HSET":    HsetHandler,
	"HGET":    HgetHandler,
	"HGETALL": HgetallHandler,
	"EXPIRE":  ExpireHandler,
	"DEL":     DelHandler,
	"KEYS":    KeysHandler,
	"ZADD":    ZaddHandler,
	"ZRANGE":  ZrangeHandler,
	"TTL":     TTLHandler,
}

var KeyValueStore = make(map[string]string)
var KeyValueStoreLock = sync.RWMutex{}

var HashStore = make(map[string]map[string]string)
var HashStoreLock = sync.RWMutex{}

var ZADDsLock sync.RWMutex
var ZADDStore = make(map[string]map[string]float64)
