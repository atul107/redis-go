package handlers

import (
	"sort"
	"strconv"

	. "github.com/redis-go/internal/aof"
	. "github.com/redis-go/pkg/resp"
)

// ZrangeHandler handles the "ZRANGE" command.
func ZrangeHandler(value Value, aof *Aof) Value {
	args := value.Array[1:]

	if len(args) < 3 {
		return Value{Typ: "error", Str: "ERR wrong number of arguments for 'zrange' command"}
	}

	setName := args[0].Bulk

	var byScore bool = true
	if len(args) > 3 {
		if args[3].Bulk == "BYSCORE" {
			byScore = true
		} else if args[3].Bulk == "BYLEX" {
			byScore = false
		}
	}

	rev := false
	if len(args) > 4 && args[4].Bulk == "REV" {
		rev = true
	}

	SortedSetStoreLock.RLock()
	set, ok := SortedSetStore[setName]
	SortedSetStoreLock.RUnlock()

	if !ok {
		return Value{Typ: "array", Array: []Value{}}
	}

	startIndexStr := args[1].Bulk
	endIndexStr := args[2].Bulk

	start, err := strconv.Atoi(startIndexStr)
	if err != nil {
		return Value{Typ: "error", Str: "ERR invalid LIMIT offset"}
	}

	end, err := strconv.Atoi(endIndexStr)
	if err != nil {
		return Value{Typ: "error", Str: "ERR invalid LIMIT count"}
	}

	if end < 0 {
		end = len(set.members) + end
	}
	if end > len(set.members)-1 {
		end = len(set.members) - 1
	}
	memberScores := make([]struct {
		member string
		score  float64
	}, end-start+1)

	for i := start; i <= end; i++ {
		memberScores[i] = struct {
			member string
			score  float64
		}{set.members[i], set.scores[i]}
	}

	if byScore {
		sortByScore(memberScores, rev)
	} else {
		sortByMember(memberScores, rev)
	}

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

		if offset < 0 {
			offset = len(memberScores) + offset
		}
		if offset < 0 {
			offset = 0
		}

		if offset >= len(memberScores) || count <= 0 {
			return Value{Typ: "array", Array: []Value{}}
		}

		endIndex := offset + count
		if endIndex > len(memberScores) {
			endIndex = len(memberScores)
		}

		memberScores = memberScores[offset:endIndex]
	}

	members := make([]Value, len(memberScores))
	for i, ms := range memberScores {
		if byScore {
			memberValue := Value{Typ: "bulk", Bulk: ms.member}
			if rev {
				members[i] = Value{Typ: "bulk", Bulk: strconv.FormatFloat(ms.score, 'f', -1, 64)}
			} else {
				members[i] = memberValue
			}
		} else {
			members[i] = Value{Typ: "bulk", Bulk: ms.member}
		}
	}

	return Value{Typ: "array", Array: members}
}

func sortByScore(memberScores []struct {
	member string
	score  float64
}, reverse bool) {
	if reverse {
		sort.SliceStable(memberScores, func(i, j int) bool {
			return memberScores[i].score > memberScores[j].score
		})
	} else {
		sort.SliceStable(memberScores, func(i, j int) bool {
			return memberScores[i].score < memberScores[j].score
		})
	}
}

func sortByMember(memberScores []struct {
	member string
	score  float64
}, reverse bool) {
	if reverse {
		sort.SliceStable(memberScores, func(i, j int) bool {
			return memberScores[i].member > memberScores[j].member
		})
	} else {
		sort.SliceStable(memberScores, func(i, j int) bool {
			return memberScores[i].member < memberScores[j].member
		})
	}
}
