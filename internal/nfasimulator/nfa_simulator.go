package nfasimulator

import (
	"unicode/utf8"

	"github.com/mmarchesotti/build-your-own-grep/internal/buildnfa"
	"github.com/mmarchesotti/build-your-own-grep/internal/nfa"
)

type capture struct {
	start int
	end   int
}

type simulationState struct {
	state     nfa.State
	lineIndex int
	captures  []capture
}

type visitedKey struct {
	state     nfa.State
	lineIndex int
}

func Simulate(line []byte, inputPattern string) (bool, [][]byte, error) {
	fragment, captureCount, err := buildnfa.Build(inputPattern)
	if err != nil {
		return false, nil, err
	}

	totalCaptures := captureCount + 1

	var statesList []simulationState
	for lineIndex := 0; lineIndex <= len(line); {
		initialCaptures := make([]capture, totalCaptures)
		for i := range initialCaptures {
			initialCaptures[i] = capture{start: -1, end: -1}
		}
		initialCaptures[0].start = lineIndex

		startState := simulationState{
			state:     fragment.Start,
			lineIndex: lineIndex,
			captures:  initialCaptures,
		}
		statesList = append(statesList, startState)

		_, size := utf8.DecodeRune(line[lineIndex:])
		if size == 0 {
			break
		}
		lineIndex += size
	}

	visited := make(map[visitedKey]bool)
	for statesIndex := 0; statesIndex < len(statesList); statesIndex++ {
		s := statesList[statesIndex]
		key := visitedKey{state: s.state, lineIndex: s.lineIndex}
		if visited[key] {
			continue
		} else {
			visited[key] = true
		}

		switch st := s.state.(type) {
		case *nfa.MatcherState:
			if s.lineIndex >= len(line) {
				break
			}

			nextRune, size := utf8.DecodeRune(line[s.lineIndex:])
			m, err := st.Matcher.Match(nextRune)
			if err != nil {
				return false, nil, err
			}
			if m {
				nextState := simulationState{
					state:     st.Out,
					lineIndex: s.lineIndex + size,
					captures:  s.captures,
				}
				statesList = append(statesList, nextState)
			}
		case *nfa.SplitState:
			nextState1 := simulationState{
				state:     st.Branch1,
				lineIndex: s.lineIndex,
				captures:  s.captures,
			}
			nextState2 := simulationState{
				state:     st.Branch2,
				lineIndex: s.lineIndex,
				captures:  s.captures,
			}
			statesList = append(statesList, nextState1)
			statesList = append(statesList, nextState2)
		case *nfa.AcceptingState:
			s.captures[0].end = s.lineIndex

			result := make([][]byte, totalCaptures)
			for i, c := range s.captures {
				if c.start != -1 && c.end != -1 {
					result[i] = line[c.start:c.end]
				}
			}
			return true, result, nil
		case *nfa.StartAnchorState:
			if s.lineIndex == 0 {
				nextState := simulationState{
					state:     st.Out,
					lineIndex: s.lineIndex,
					captures:  s.captures,
				}
				statesList = append(statesList, nextState)
			}
		case *nfa.EndAnchorState:
			if s.lineIndex == len(line) {
				nextState := simulationState{
					state:     st.Out,
					lineIndex: s.lineIndex,
					captures:  s.captures,
				}
				statesList = append(statesList, nextState)
			}
		case *nfa.StartCaptureState:
			newCaptures := make([]capture, len(s.captures))
			copy(newCaptures, s.captures)
			newCaptures[st.Index].start = s.lineIndex

			nextState := simulationState{
				state:     st.Out,
				lineIndex: s.lineIndex,
				captures:  newCaptures,
			}
			statesList = append(statesList, nextState)

		case *nfa.EndCaptureState:
			newCaptures := make([]capture, len(s.captures))
			copy(newCaptures, s.captures)
			newCaptures[st.Index].end = s.lineIndex

			nextState := simulationState{
				state:     st.Out,
				lineIndex: s.lineIndex,
				captures:  newCaptures,
			}
			statesList = append(statesList, nextState)
		}
	}

	return false, nil, nil
}
