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

	SortedSetStoreLock.RLock()
	set, ok := SortedSetStore[setName]
	SortedSetStoreLock.RUnlock()

	if !ok {
		return Value{Typ: "array", Array: []Value{}}
	}

	startIndexStr := args[1].Bulk
	endIndexStr := args[2].Bulk

	excludeLeft := false
	if startIndexStr[0] == '(' {
		excludeLeft = true
		startIndexStr = startIndexStr[1:]
	}

	excludeRight := false
	if endIndexStr[0] == '(' {
		excludeRight = true
		endIndexStr = endIndexStr[1:]
	}

	start, err := strconv.Atoi(startIndexStr)
	if err != nil {
		return Value{Typ: "error", Str: "ERR invalid LIMIT offset"}
	}

	if start < 0 {
		start = len(set.members) + start
	}
	if start < 0 {
		start = 0
	}

	if excludeLeft {
		start = start + 1
	}
	end, err := strconv.Atoi(endIndexStr)
	if err != nil {
		return Value{Typ: "error", Str: "ERR invalid LIMIT count"}
	}

	if end < 0 {
		end = len(set.members) + end
	}
	if end < 0 {
		end = 0
	}
	if end > len(set.members)-1 {
		end = len(set.members) - 1
	}

	if excludeRight {
		end = end - 1
	}

	memberScores := make([]struct {
		member string
		score  float64
	}, end-start+1)

	memberScoresIdx := 0
	for i := start; i <= end; i++ {
		memberScores[memberScoresIdx] = struct {
			member string
			score  float64
		}{set.members[i], set.scores[i]}
		memberScoresIdx++
	}

	var byScore, rev, withScores, limitReq bool
	var offsetStr, countStr string
	byScore = true
	for i := 3; i < len(args); i++ {
		arg := args[i].Bulk

		if arg == "BYSCORE" {
			byScore = true
		} else if arg == "BYLEX" {
			byScore = false
		} else if arg == "REV" {
			rev = true
		} else if arg == "WITHSCORES" {
			withScores = true
		} else if arg == "LIMIT" {
			limitReq = true
			offsetStr = args[i+1].Bulk
			i++
			countStr = args[i+1].Bulk
			i++
		}
	}

	if byScore {
		sortByScore(memberScores, rev)
	} else {
		sortByMember(memberScores, rev)
	}

	if limitReq {
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

	if withScores {
		membersWithScores := make([]Value, len(memberScores)*2)
		for i, ms := range memberScores {
			// if byScore {
			membersWithScores[i*2] = Value{Typ: "bulk", Bulk: ms.member}
			membersWithScores[i*2+1] = Value{Typ: "bulk", Bulk: strconv.FormatFloat(ms.score, 'f', -1, 64)}
			// } else {
			// 	membersWithScores[i*2] = Value{Typ: "bulk", Bulk: ms.member}
			// }
		}
		return Value{Typ: "array", Array: membersWithScores}
	} else {
		members := make([]Value, len(memberScores))
		for i, ms := range memberScores {
			// if byScore {
			// 	memberValue := Value{Typ: "bulk", Bulk: ms.member}
			// 	if rev {
			// 		members[i] = Value{Typ: "bulk", Bulk: strconv.FormatFloat(ms.score, 'f', -1, 64)}
			// 	} else {
			// 		members[i] = memberValue
			// 	}
			// } else {
			members[i] = Value{Typ: "bulk", Bulk: ms.member}
			// }
		}
		return Value{Typ: "array", Array: members}
	}
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
