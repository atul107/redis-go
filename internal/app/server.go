package app

import (
	"fmt"
	"net"
	"strings"
	"time"

	. "github.com/redis-go/config"
	. "github.com/redis-go/internal/aof"
	. "github.com/redis-go/internal/handlers"
	. "github.com/redis-go/pkg/resp"
)

func Run(config *Config) {
	fmt.Println("Listening on port: ", config.TCP.Addr)

	// Create a new server
	l, err := net.Listen("tcp", config.TCP.Addr)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer l.Close()

	aof, err := NewAof(config.AoF.FileName)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer aof.Close()

	go expirationCheck(aof)

	// aof.Read(KeyValueStore, ZADDStore, func(value Value) {
	// 	command := strings.ToUpper(value.Array[0].Bulk)

	// 	handler, ok := Handlers[command]
	// 	if !ok {
	// 		fmt.Println("Invalid command: ", command)
	// 		return
	// 	}

	// 	handler(value, aof)
	// })

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println(err)
			return
		}
		go handleConnection(conn, aof)
	}
}

func handleConnection(conn net.Conn, aof *Aof) {
	defer conn.Close()

	for {
		resp := NewRespReader(conn)
		value, err := resp.Read()
		if err != nil {
			fmt.Println(err)
			return
		}

		writer := NewRespWriter(conn)

		if value.Typ != "array" {
			fmt.Println("Invalid request, expected array")
			writer.WriteValue(Value{Typ: "string", Str: "Invalid request"})
			continue
		}

		if len(value.Array) == 0 {
			fmt.Println("Invalid request, expected array length > 0")
			writer.WriteValue(Value{Typ: "string", Str: "Invalid request"})
			continue
		}

		command := strings.ToUpper(value.Array[0].Bulk)

		handler, ok := Handlers[command]
		if !ok {
			fmt.Println("Invalid command: ", command)
			writer.WriteValue(Value{Typ: "string", Str: ""})
			continue
		}

		result := handler(value, aof)
		writer.WriteValue(result)
	}
}

func expirationCheck(aof *Aof) {
	for {
		time.Sleep(1 * time.Second)

		currentTimestamp := time.Now()

		ExpirySToreLock.Lock()
		for key, expireTime := range ExpiryStore {
			if currentTimestamp.After(expireTime) {
				KeyValueStoreLock.Lock()
				delete(KeyValueStore, key)
				KeyValueStoreLock.Unlock()

				delete(ExpiryStore, key)

				aof.Write(Value{Typ: "array", Array: []Value{
					{Typ: "bulk", Bulk: "DEL"},
					{Typ: "bulk", Bulk: key},
				}})
			}
		}
		ExpirySToreLock.Unlock()
	}
}
