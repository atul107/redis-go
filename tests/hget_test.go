package tests

import "testing"

func TestHgetHandler(t *testing.T) {
	t.Run("HGET on Existing Hash and Key", func(t *testing.T) {
		conn, cleanup := connectToServer(t)
		defer cleanup()

		sendCommand(t, conn, "HSET myhash mykey myvalue")

		response := sendCommand(t, conn, "HGET myhash mykey")
		assertResponse(t, response, "myvalue")
	})

	t.Run("HGET on Existing Hash and Non-Existent Key", func(t *testing.T) {
		conn, cleanup := connectToServer(t)
		defer cleanup()

		sendCommand(t, conn, "HSET myhash mykey myvalue")

		response := sendCommand(t, conn, "HGET myhash nonkey")
		assertResponse(t, response, "(nil)")
	})

	t.Run("HGET with Wrong Number of Arguments", func(t *testing.T) {
		conn, cleanup := connectToServer(t)
		defer cleanup()

		response := sendCommand(t, conn, "HGET myhash")
		assertResponse(t, response, "ERR wrong number of arguments for 'hget' command")
	})

	t.Log("All GET handler test cases passed.")
}
