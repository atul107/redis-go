package handlers

import (
	"fmt"
	"strconv"

	. "github.com/redis-go/internal/aof"
	. "github.com/redis-go/pkg/resp"
)

type SortedSet struct {
	members []string
	scores  []float64
}

// ZaddHandler handles the "ZADD" command.
func ZaddHandler(value Value, aof *Aof) Value {
	args := value.Array[1:]

	if len(args) < 3 {
		return Value{Typ: "error", Str: "ERR wrong number of arguments for 'zadd' command"}
	}

	setName := args[0].Bulk
	var ch bool
	var nx, xx, gt, lt, incr bool

	SortedSetStoreLock.Lock()
	set, ok := SortedSetStore[setName]
	if !ok {
		set = &SortedSet{members: []string{}, scores: []float64{}}
		SortedSetStore[setName] = set
	}
	SortedSetStoreLock.Unlock()

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

				SortedSetStoreLock.Lock()
				index := -1
				for j, m := range set.members {
					if m == member {
						index = j
						break
					}
				}

				if (nx && index >= 0) || (xx && index < 0) {
					SortedSetStoreLock.Unlock()
					return Value{Typ: "integer", Num: 0}
				}

				if incr {
					if index < 0 {
						SortedSetStoreLock.Unlock()
						return Value{Typ: "error", Str: "ERR INCR option specified, but member does not exist"}
					}
					score += set.scores[index]
				}

				if (gt && lt) || (gt && score == 0) || (lt && score == 0) {
					SortedSetStoreLock.Unlock()
					return Value{Typ: "error", Str: "ERR syntax error"}
				}

				if gt || lt {
					newMembers := []string{}
					newScores := []float64{}
					for j, s := range set.scores {
						if (gt && s > score) || (lt && s < score) {
							newMembers = append(newMembers, set.members[j])
							newScores = append(newScores, s)
						}
					}
					set.members = newMembers
					set.scores = newScores
				}

				if index >= 0 {
					set.scores[index] = score
				} else {
					set.members = append(set.members, member)
					set.scores = append(set.scores, score)
				}

				for j := len(set.members) - 1; j > 0; j-- {
					if set.scores[j] < set.scores[j-1] {
						set.members[j], set.members[j-1] = set.members[j-1], set.members[j]
						set.scores[j], set.scores[j-1] = set.scores[j-1], set.scores[j]
					}
				}

				SortedSetStoreLock.Unlock()

				if !ch || (ch && index < 0) {
					aof.Write(Value{Typ: "array", Array: []Value{
						{Typ: "bulk", Bulk: "ZADD"},
						{Typ: "bulk", Bulk: setName},
						{Typ: "bulk", Bulk: fmt.Sprintf("%f", score)},
						{Typ: "bulk", Bulk: member},
					}})
				}
			}
			i++
		}
	}

	return Value{Typ: "integer", Num: len(args) / 2}
}
