package resp

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

// RespReader represents a RESP value reader.
type RespReader struct {
	reader *bufio.Reader
}

// NewRespReader creates a new RESP reader.
func NewRespReader(r io.Reader) *RespReader {
	return &RespReader{reader: bufio.NewReader(r)}
}

func (r *RespReader) readLine() (line []byte, n int, err error) {
	for {
		b, err := r.reader.ReadByte()
		if err != nil {
			return nil, 0, err
		}
		n += 1
		line = append(line, b)
		if len(line) >= 2 && line[len(line)-2] == '\r' {
			break
		}
	}
	return line[:len(line)-2], n, nil
}

func (r *RespReader) readInteger() (x int, n int, err error) {
	line, n, err := r.readLine()
	if err != nil {
		return 0, 0, err
	}
	i64, err := strconv.ParseInt(string(line), 10, 64)
	if err != nil {
		return 0, n, err
	}
	return int(i64), n, nil
}

func (r *RespReader) Read() (Value, error) {
	_type, err := r.reader.ReadByte()

	if err != nil {
		return Value{}, err
	}

	switch _type {
	case ARRAY:
		return r.readArray()
	case BULK:
		return r.readBulk()
	default:
		fmt.Printf("Unknown type: %v", string(_type))
		return Value{}, nil
	}
}

func (r *RespReader) readArray() (Value, error) {
	v := Value{}
	v.Typ = "array"

	// read length of array
	len, _, err := r.readInteger()
	if err != nil {
		return v, err
	}

	// foreach line, parse and read the value
	v.Array = make([]Value, 0)
	for i := 0; i < len; i++ {
		val, err := r.Read()
		if err != nil {
			return v, err
		}

		// append parsed value to array
		v.Array = append(v.Array, val)
	}

	return v, nil
}

func (r *RespReader) readBulk() (Value, error) {
	v := Value{}

	v.Typ = "bulk"

	len, _, err := r.readInteger()
	if err != nil {
		return v, err
	}

	bulk := make([]byte, len)

	r.reader.Read(bulk)

	v.Bulk = string(bulk)

	// Read the trailing CRLF
	r.readLine()

	return v, nil
}

func (r *RespReader) SendCommand(writer *bufio.Writer, args ...string) {
	// Format the command as RESP and send it to the server
	command := "*" + fmt.Sprint(len(args)) + "\r\n"
	for _, arg := range args {
		command += "$" + fmt.Sprint(len(arg)) + "\r\n" + arg + "\r\n"
	}
	writer.WriteString(command)
	writer.Flush()
}

func (r *RespReader) ReadResponse() (string, error) {
	// Read the first byte to determine the response type
	responseType, err := r.reader.ReadByte()
	if err != nil {
		return "", err
	}
	// Read the response based on its type
	switch responseType {
	case '+':
		// Simple string
		response, err := r.reader.ReadString('\n')
		if err != nil {
			return "", err
		}
		return strings.TrimSpace(response), nil
	case '-':
		// Error
		response, err := r.reader.ReadString('\n')
		if err != nil {
			return "", err
		}
		return strings.TrimSpace(response), nil
	case ':':
		// Integer
		response, err := r.reader.ReadString('\n')
		if err != nil {
			return "", err
		}
		return strings.TrimSpace(response), nil
	case '$':
		// Bulk string
		length, err := r.reader.ReadString('\n')
		if err != nil {
			return "", err
		}
		var stringLength int
		if _, err := fmt.Sscanf(strings.TrimSpace(length), "%d", &stringLength); err != nil {
			return "", fmt.Errorf("Invalid bulk string length: %v, %v, %v, %v", length, stringLength, responseType, err)
		}

		if stringLength == -1 {
			// Null bulk string
			return "(nil)", nil
		}

		data := make([]byte, stringLength)
		_, err = io.ReadFull(r.reader, data)
		if err != nil {
			return "", err
		}
		// Read and discard the trailing '\r\n'
		_, _ = r.reader.Discard(2)
		return string(data), nil
	case '*':
		// Array
		arrayLength, err := r.reader.ReadString('\n')
		if err != nil {
			return "", err
		}
		var arraySize int
		if _, err := fmt.Sscanf(strings.TrimSpace(arrayLength), "%d", &arraySize); err != nil {
			return "", fmt.Errorf("Invalid array length")
		}
		// Read the elements of the array
		var result []string
		for i := 0; i < arraySize; i++ {
			element, err := r.ReadResponse()
			if err != nil {
				return "", err
			}
			result = append(result, element)
		}
		return strings.Join(result, ", "), nil
	default:
		return "", fmt.Errorf("Invalid response type: %c", responseType)
	}
}

func (r *RespReader) ReadInput() string {
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}
