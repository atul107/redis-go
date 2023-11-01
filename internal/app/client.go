package app

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"

	. "github.com/redis-go/config"
	. "github.com/redis-go/pkg/resp"
)

func RunClient(config *Config) {
	serverAddr := config.TCP.Addr

	// Connect to the Redis server
	conn, err := net.Dial("tcp", serverAddr)
	if err != nil {
		fmt.Println("Error connecting to Redis server:", err)
		os.Exit(1)
	}
	defer conn.Close()

	// Create a reader and writer for the connection
	reader := NewRespReader(conn)
	writer := bufio.NewWriter(conn)

	for {
		fmt.Print(config.TCP.Addr, "> ")
		input := reader.ReadInput()

		if strings.ToLower(input) == "quit" {
			break
		}

		// Send the user input as a Redis command
		reader.SendCommand(writer, strings.Split(input, " ")...)

		response, err := reader.ReadResponse()
		if err != nil {
			fmt.Println("Error reading response:", err)
		} else {
			fmt.Println(response)
		}
	}
}
