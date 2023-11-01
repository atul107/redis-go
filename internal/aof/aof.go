package aof

import (
	"bufio"
	"io"
	"os"
	"strconv"
	"sync"

	. "github.com/redis-go/pkg/resp"
)

type Aof struct {
	file *os.File
	rd   *bufio.Reader
	mu   sync.Mutex
}

func NewAof(path string) (*Aof, error) {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		return nil, err
	}

	aof := &Aof{
		file: f,
		rd:   bufio.NewReader(f),
	}

	// Start a goroutine to sync AOF to disk every 1 second
	// go func() {
	// 	for {
	// 		aof.mu.Lock()
	// 		aof.file.Sync()
	// 		aof.mu.Unlock()
	// 		time.Sleep(time.Second)
	// 	}
	// }()

	return aof, nil
}

func (aof *Aof) Close() error {
	aof.mu.Lock()
	defer aof.mu.Unlock()
	return aof.file.Close()
}

func (aof *Aof) Write(value Value) error {
	aof.mu.Lock()
	defer aof.mu.Unlock()

	_, err := aof.file.Write(value.Marshal())
	if err != nil {
		return err
	}

	return nil
}

func (aof *Aof) Read(keyValueStore map[string]string, ZADDStore map[string]map[string]float64, fn func(value Value)) error {
	aof.mu.Lock()
	defer aof.mu.Unlock()

	_, err := aof.file.Seek(0, io.SeekStart)
	if err != nil {
		return err
	}

	reader := NewRespReader(aof.file)

	for {
		value, err := reader.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}

		// Check for DEL commands and execute them
		if value.Typ == "array" && len(value.Array) >= 2 &&
			value.Array[0].Bulk == "DEL" {
			key := value.Array[1].Bulk
			// keyValueStoreLock.Lock()
			delete(keyValueStore, key)
			// keyValueStoreLock.Unlock()
		}

		if value.Typ == "array" && len(value.Array) >= 4 &&
			value.Array[0].Bulk == "ZADD" {
			setName := value.Array[1].Bulk
			scoreStr := value.Array[2].Bulk
			member := value.Array[3].Bulk

			score, err := strconv.ParseFloat(scoreStr, 64)
			if err != nil {
				// Handle error
			}

			// Add the member and score to the sorted set
			// ZADDsLock.Lock()
			if _, ok := ZADDStore[setName]; !ok {
				ZADDStore[setName] = map[string]float64{}
			}
			ZADDStore[setName][member] = score
			// ZADDsLock.Unlock()
		}

		fn(value)
	}

	return nil
}
