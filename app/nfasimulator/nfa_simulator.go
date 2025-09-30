package nfasimulator

import (
	"unicode/utf8"

	"github.com/mmarchesotti/build-your-own-grep/app/buildnfa"
	"github.com/mmarchesotti/build-your-own-grep/app/nfa"
)

type simulationState struct {
	state     nfa.State
	lineIndex int
}

func Simulate(line []byte, inputPattern string) (bool, error) {
	fragment, err := buildnfa.Build(inputPattern)
	if err != nil {
		return false, err
	}
	var statesList []simulationState
	for lineIndex := 0; lineIndex <= len(line); {
		startState := simulationState{
			state:     fragment.Start,
			lineIndex: lineIndex,
		}
		statesList = append(statesList, startState)

		_, size := utf8.DecodeRune(line[lineIndex:])
		if size == 0 {
			break
		}
		lineIndex += size
	}

	visited := make(map[simulationState]bool)
	for statesIndex := 0; statesIndex < len(statesList); statesIndex++ {
		s := statesList[statesIndex]
		if visited[s] {
			continue
		} else {
			visited[s] = true
		}
		switch st := s.state.(type) {
		case *nfa.MatcherState:
			if s.lineIndex >= len(line) {
				break
			}

			nextRune, size := utf8.DecodeRune(line[s.lineIndex:])
			m, err := st.Matcher.Match(nextRune)
			if err != nil {
				return false, err
			}
			if m {
				nextState := simulationState{
					state:     st.Out,
					lineIndex: s.lineIndex + size,
				}
				statesList = append(statesList, nextState)
			}
		case *nfa.SplitState:
			nextState1 := simulationState{
				state:     st.Branch1,
				lineIndex: s.lineIndex,
			}
			nextState2 := simulationState{
				state:     st.Branch2,
				lineIndex: s.lineIndex,
			}
			statesList = append(statesList, nextState1)
			statesList = append(statesList, nextState2)
		case *nfa.AcceptingState:
			return true, nil
		case *nfa.StartAnchorState:
			if s.lineIndex == 0 {
				nextState := simulationState{
					state:     st.Out,
					lineIndex: s.lineIndex,
				}
				statesList = append(statesList, nextState)
			}
		case *nfa.EndAnchorState:
			if s.lineIndex == len(line) {
				nextState := simulationState{
					state:     st.Out,
					lineIndex: s.lineIndex,
				}
				statesList = append(statesList, nextState)
			}
		}
	}

	return false, nil
}
