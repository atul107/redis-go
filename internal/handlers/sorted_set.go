package handlers

import (
	"fmt"
	"strconv"
	"sync"

	. "github.com/redis-go/internal/aof"
	. "github.com/redis-go/pkg/resp"
)

var SortedSets = map[string]map[string]float64{}
var SortedSetsMu = sync.RWMutex{}

func ZaddHandler(value Value, aof *Aof) Value {
	args := value.Array[1:]

	if len(args) < 3 {
		return Value{Typ: "error", Str: "ERR wrong number of arguments for 'zadd' command"}
	}

	setName := args[0].Bulk

	// Parse the member-score pairs and add them to the sorted set
	for i := 1; i < len(args); i += 2 {
		member := args[i].Bulk
		scoreStr := args[i+1].Bulk
		score, err := strconv.ParseFloat(scoreStr, 64)
		if err != nil {
			return Value{Typ: "error", Str: "ERR invalid score"}
		}

		// Add the member and score to the sorted set
		ZADDsLock.Lock()
		if _, ok := ZADDStore[setName]; !ok {
			ZADDStore[setName] = map[string]float64{}
		}
		ZADDStore[setName][member] = score
		ZADDsLock.Unlock()

		// Record the ZADD command in the AOF
		aof.Write(Value{Typ: "array", Array: []Value{
			{Typ: "bulk", Bulk: "ZADD"},
			{Typ: "bulk", Bulk: setName},
			{Typ: "bulk", Bulk: fmt.Sprintf("%f", score)},
			{Typ: "bulk", Bulk: member},
		}})
	}

	return Value{Typ: "integer", Num: len(args) / 2}
}

func ZrangeHandler(value Value, aof *Aof) Value {
	args := value.Array[1:]

	if len(args) < 3 {
		return Value{Typ: "error", Str: "ERR wrong number of arguments for 'zrange' command"}
	}

	setName := args[0].Bulk
	startStr := args[1].Bulk
	stopStr := args[2].Bulk

	start, err := strconv.Atoi(startStr)
	if err != nil {
		return Value{Typ: "error", Str: "ERR invalid start index"}
	}

	stop, err := strconv.Atoi(stopStr)
	if err != nil {
		return Value{Typ: "error", Str: "ERR invalid stop index"}
	}

	// Ensure the sorted set exists
	ZADDsLock.RLock()
	sortedSet, ok := ZADDStore[setName]
	ZADDsLock.RUnlock()

	if !ok {
		return Value{Typ: "array", Array: []Value{}}
	}

	// Retrieve members in the specified range
	members := []Value{}
	index := 0

	for member, score := range sortedSet {
		if index >= start && index <= stop {
			members = append(members, Value{Typ: "bulk", Bulk: member})
			members = append(members, Value{Typ: "bulk", Bulk: strconv.FormatFloat(score, 'f', -1, 64)})
		}
		index++
	}

	return Value{Typ: "array", Array: members}
}
