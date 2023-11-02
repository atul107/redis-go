package handlers

import (
	"fmt"
	"strconv"

	. "github.com/redis-go/internal/aof"
	. "github.com/redis-go/pkg/resp"
)

// ZaddHandler handles the "ZADD" command.
func ZaddHandler(value Value, aof *Aof) Value {
	args := value.Array[1:]

	if len(args) < 3 {
		return Value{Typ: "error", Str: "ERR wrong number of arguments for 'zadd' command"}
	}

	setName := args[0].Bulk
	var memberExists bool
	var ch bool
	var nx, xx, gt, lt, incr bool

	for i := 1; i < len(args); i++ {
		arg := args[i].Bulk

		if arg == "NX" {
			nx = true
		} else if arg == "XX" {
			xx = true
		} else if arg == "GT" {
			gt = true
		} else if arg == "LT" {
			lt = true
		} else if arg == "CH" {
			ch = true
		} else if arg == "INCR" {
			incr = true
		} else {
			if i+1 < len(args) {
				scoreStr := arg
				member := args[i+1].Bulk
				score, err := strconv.ParseFloat(scoreStr, 64)
				if err != nil {
					return Value{Typ: "error", Str: "ERR invalid score"}
				}

				ZADDStoreLock.Lock()
				if set, ok := ZADDStore[setName]; ok {
					_, memberExists = set[member]

					if (nx && memberExists) || (xx && !memberExists) {
						ZADDStoreLock.Unlock()
						return Value{Typ: "integer", Num: 0}
					}

					if incr {
						if !memberExists {
							ZADDStoreLock.Unlock()
							return Value{Typ: "error", Str: "ERR INCR option specified, but member does not exist"}
						}
						score += set[member]
					}
				}

				if (gt && lt) || (gt && score == 0) || (lt && score == 0) {
					ZADDStoreLock.Unlock()
					return Value{Typ: "error", Str: "ERR syntax error"}
				}

				if gt {
					newSet := make(map[string]float64)
					for k, v := range ZADDStore[setName] {
						if v > score {
							newSet[k] = v
						}
					}
					ZADDStore[setName] = newSet
				}

				if lt {
					newSet := make(map[string]float64)
					for k, v := range ZADDStore[setName] {
						if v < score {
							newSet[k] = v
						}
					}
					ZADDStore[setName] = newSet
				}

				if !gt && !lt {
					ZADDStore[setName] = make(map[string]float64)
				}

				ZADDStore[setName][member] = score
				ZADDStoreLock.Unlock()

				if !ch || (ch && !memberExists) {
					aof.Write(Value{Typ: "array", Array: []Value{
						{Typ: "bulk", Bulk: "ZADD"},
						{Typ: "bulk", Bulk: setName},
						{Typ: "bulk", Bulk: fmt.Sprintf("%f", score)},
						{Typ: "bulk", Bulk: member},
					}})
				}
				i++
			}
		}
	}

	return Value{Typ: "integer", Num: len(args) / 2}
}
