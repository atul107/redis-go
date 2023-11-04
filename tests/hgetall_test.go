package tests

import "testing"

func TestHgetallHandler(t *testing.T) {
	t.Run("HGETALL on Existing Hash", func(t *testing.T) {
		conn, cleanup := connectToServer(t)
		defer cleanup()

		sendCommand(t, conn, "HSET hashkey key1 value1")
		sendCommand(t, conn, "HSET hashkey key2 value2")

		response := sendCommand(t, conn, "HGETALL hashkey")

		expectedResponse := "key1, value1, key2, value2"
		assertResponse(t, response, expectedResponse)
	})

	t.Run("HGETALL on Non-Existent Hash", func(t *testing.T) {
		conn, cleanup := connectToServer(t)
		defer cleanup()

		response := sendCommand(t, conn, "HGETALL newhashkey")

		assertResponse(t, response, "(nil)")
	})

	t.Run("HGETALL with Wrong Number of Arguments", func(t *testing.T) {
		conn, cleanup := connectToServer(t)
		defer cleanup()

		response := sendCommand(t, conn, "HGETALL myhash key1")
		assertResponse(t, response, "ERR wrong number of arguments for 'hgetall' command")
	})

	t.Log("All HGETALL handler test cases passed.")
}
