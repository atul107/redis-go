package tests

import "testing"

func TestHsetHandler(t *testing.T) {
	t.Run("HSET on Existing Hash and Key", func(t *testing.T) {
		conn, cleanup := connectToServer(t)
		defer cleanup()

		response := sendCommand(t, conn, "HSET myhash mykey myvalue")
		assertResponse(t, response, "OK")

		response = sendCommand(t, conn, "HGET myhash mykey")
		assertResponse(t, response, "myvalue")
	})

	t.Run("HSET on Existing Hash and New Key", func(t *testing.T) {
		conn, cleanup := connectToServer(t)
		defer cleanup()

		response := sendCommand(t, conn, "HSET myhash mykey myvalue")
		assertResponse(t, response, "OK")

		response = sendCommand(t, conn, "HSET myhash newkey newvalue")
		assertResponse(t, response, "OK")

		response = sendCommand(t, conn, "HGET myhash newkey")
		assertResponse(t, response, "newvalue")
	})

	t.Run("HSET on New Hash", func(t *testing.T) {
		conn, cleanup := connectToServer(t)
		defer cleanup()

		response := sendCommand(t, conn, "HSET newhash key1 value1")
		assertResponse(t, response, "OK")

		response = sendCommand(t, conn, "HGET newhash key1")
		assertResponse(t, response, "value1")
	})

	t.Run("HSET with Wrong Number of Arguments", func(t *testing.T) {
		conn, cleanup := connectToServer(t)
		defer cleanup()

		response := sendCommand(t, conn, "HSET myhash mykey")
		assertResponse(t, response, "ERR wrong number of arguments for 'hset' command")
	})

	t.Log("All HSET handler test cases passed.")
}
