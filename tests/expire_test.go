package tests

import "testing"

func TestExpireHandler(t *testing.T) {
	t.Run("Set Expiration with Valid Key and Time", func(t *testing.T) {
		conn, cleanup := connectToServer(t)
		defer cleanup()

		response := sendCommand(t, conn, "SET mykey myvalue")
		assertResponse(t, response, "OK")

		response = sendCommand(t, conn, "EXPIRE mykey 10")
		assertResponse(t, response, "1")
	})

	t.Run("Set Expiration with Key Already Having Expiration", func(t *testing.T) {
		conn, cleanup := connectToServer(t)
		defer cleanup()

		response := sendCommand(t, conn, "SET mykey myvalue")
		assertResponse(t, response, "OK")

		response = sendCommand(t, conn, "EXPIRE mykey 10")
		assertResponse(t, response, "1")

		response = sendCommand(t, conn, "EXPIRE mykey 20")
		assertResponse(t, response, "1")
	})

	t.Run("Set Expiration with NX Option", func(t *testing.T) {
		conn, cleanup := connectToServer(t)
		defer cleanup()

		response := sendCommand(t, conn, "SET mykey myvalue")
		assertResponse(t, response, "OK")

		response = sendCommand(t, conn, "EXPIRE mykey 20 NX")
		assertResponse(t, response, "0")
	})

	t.Run("Set Expiration with XX Option", func(t *testing.T) {
		conn, cleanup := connectToServer(t)
		defer cleanup()

		response := sendCommand(t, conn, "EXPIRE nonexisting_key 20 XX")
		assertResponse(t, response, "0")

		response = sendCommand(t, conn, "SET mykey myvalue")
		assertResponse(t, response, "OK")

		response = sendCommand(t, conn, "EXPIRE mykey 20 XX")
		assertResponse(t, response, "1")
	})

	t.Run("Set Expiration with GT Option", func(t *testing.T) {
		conn, cleanup := connectToServer(t)
		defer cleanup()

		response := sendCommand(t, conn, "SET mykey myvalue")
		assertResponse(t, response, "OK")

		response = sendCommand(t, conn, "EXPIRE mykey 100")
		assertResponse(t, response, "1")

		response = sendCommand(t, conn, "EXPIRE mykey 10 GT")
		assertResponse(t, response, "0")

		response = sendCommand(t, conn, "EXPIRE mykey 200 GT")
		assertResponse(t, response, "1")
	})

	t.Run("Set Expiration with LT Option", func(t *testing.T) {
		conn, cleanup := connectToServer(t)
		defer cleanup()

		response := sendCommand(t, conn, "SET mykey myvalue")
		assertResponse(t, response, "OK")

		response = sendCommand(t, conn, "EXPIRE mykey 100")
		assertResponse(t, response, "1")

		response = sendCommand(t, conn, "EXPIRE mykey 10 LT")
		assertResponse(t, response, "1")

		response = sendCommand(t, conn, "EXPIRE mykey 200 LT")
		assertResponse(t, response, "0")
	})
	t.Log("All EXPIRE handler test cases passed.")
}
