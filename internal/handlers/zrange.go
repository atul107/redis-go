package handlers

import (
	"sort"
	"strconv"
	"sync"

	. "github.com/redis-go/internal/aof"
	. "github.com/redis-go/pkg/resp"
)

var SortedSets = map[string]map[string]float64{}
var SortedSetsMu = sync.RWMutex{}

func ZrangeHandler(value Value, aof *Aof) Value {
	args := value.Array[1:]

	if len(args) < 3 {
		return Value{Typ: "error", Str: "ERR wrong number of arguments for 'zrange' command"}
	}

	setName := args[0].Bulk

	// Check if the "BYSCORE" or "BYLEX" option is specified
	var byScore bool
	if len(args) > 3 {
		if args[3].Bulk == "BYSCORE" {
			byScore = true
		} else if args[3].Bulk == "BYLEX" {
			byScore = false
		}
	}

	// Check if the "REV" option is specified
	rev := false
	if len(args) > 4 && args[4].Bulk == "REV" {
		rev = true
	}

	// Initialize the member-score pairs
	memberScores := []struct {
		member string
		score  float64
	}{}

	// Ensure the sorted set exists
	ZADDStoreLock.RLock()
	sortedSet, ok := ZADDStore[setName]
	ZADDStoreLock.RUnlock()

	if !ok {
		return Value{Typ: "array", Array: []Value{}}
	}

	// Collect member-score pairs
	for member, score := range sortedSet {
		memberScores = append(memberScores, struct {
			member string
			score  float64
		}{member, score})
	}

	// Sort the member-score pairs by score or lexicographically
	if byScore {
		sort.SliceStable(memberScores, func(i, j int) bool {
			return memberScores[i].score < memberScores[j].score
		})
	} else {
		sort.SliceStable(memberScores, func(i, j int) bool {
			return memberScores[i].member < memberScores[j].member
		})
	}

	// Check if the "LIMIT" option is specified
	if len(args) > 5 && args[5].Bulk == "LIMIT" {
		if len(args) < 8 {
			return Value{Typ: "error", Str: "ERR wrong number of arguments for 'zrange' command with 'LIMIT'"}
		}

		offsetStr := args[6].Bulk
		countStr := args[7].Bulk

		offset, err := strconv.Atoi(offsetStr)
		if err != nil {
			return Value{Typ: "error", Str: "ERR invalid LIMIT offset"}
		}

		count, err := strconv.Atoi(countStr)
		if err != nil {
			return Value{Typ: "error", Str: "ERR invalid LIMIT count"}
		}

		// Apply the offset and count
		if offset < 0 {
			offset = len(memberScores) + offset
		}
		if offset < 0 {
			offset = 0
		}

		if offset+count > len(memberScores) {
			count = len(memberScores) - offset
		}

		memberScores = memberScores[offset : offset+count]
	}

	// Prepare the response with or without scores
	members := []Value{}
	for _, ms := range memberScores {
		if byScore {
			members = append(members, Value{Typ: "bulk", Bulk: ms.member})
			if rev {
				members = append(members, Value{Typ: "bulk", Bulk: strconv.FormatFloat(ms.score, 'f', -1, 64)})
			}
		} else {
			members = append(members, Value{Typ: "bulk", Bulk: ms.member})
		}
	}

	return Value{Typ: "array", Array: members}
}
