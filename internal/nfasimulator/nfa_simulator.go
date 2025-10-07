package nfasimulator

import (
	"fmt"
	"unicode/utf8"

	"github.com/mmarchesotti/build-your-own-grep/internal/buildnfa"
	"github.com/mmarchesotti/build-your-own-grep/internal/nfa"
)

// Represents a position in the input string and the NFA state.
type thread struct {
	state     nfa.State
	lineIndex int
	captures  []int
}

// key generates a unique, comparable string key for the thread.
func (t *thread) key() string {
	// Note: %p prints the pointer address of the state.
	// This is crucial because different states can be at the same memory address
	// in different simulation runs, but we only care about the state's identity within one run.
	return fmt.Sprintf("%p-%d", t.state, t.lineIndex)
}

// Represents an action to take (either exploring a node or reverting from it).
type task struct {
	isRevert bool
	thread   thread
	undoLog  []undoEntry
}

// Stores information needed to undo a capture.
type undoEntry struct {
	captureIndex int
	oldValue     int
}

// Simulate finds the first match of the pattern in the line.
func Simulate(line []byte, inputPattern string) (bool, []int, error) {
	fragment, captureCount, err := buildnfa.Build(inputPattern)
	if err != nil {
		return false, nil, err
	}

	for i := 0; i <= len(line); i++ {
		// Try to find a match starting at each position `i` in the line.
		captures, found := findMatchAt(fragment.Start, line, i, captureCount)
		if found {
			return true, captures, nil
		}
		// If we are at the end of the line, don't advance further.
		if i == len(line) {
			break
		}
	}
	return false, nil, nil
}

// findMatchAt attempts to find a match starting from a specific index.
func findMatchAt(startState nfa.State, line []byte, startIndex int, captureCount int) ([]int, bool) {
	stack := []task{}

	// Initial captures array, -1 indicates not set.
	initialCaptures := make([]int, captureCount*2)
	for i := range initialCaptures {
		initialCaptures[i] = -1
	}

	// Start the DFS from the initial state and position.
	initialThread := thread{state: startState, lineIndex: startIndex, captures: initialCaptures}
	stack = append(stack, task{isRevert: false, thread: initialThread, undoLog: nil})

	// CORRECTED: Use a map with a string key instead of the struct itself.
	visited := make(map[string]bool)

	for len(stack) > 0 {
		// Pop the next task.
		currentTask := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		// If it's a revert task, undo captures.
		if currentTask.isRevert {
			for _, entry := range currentTask.undoLog {
				currentTask.thread.captures[entry.captureIndex] = entry.oldValue
			}
			continue
		}

		// CORRECTED: Use the thread's key() method for the visited check.
		threadKey := currentTask.thread.key()
		if visited[threadKey] {
			continue
		}
		visited[threadKey] = true

		// Process the current state.
		currentState := currentTask.thread.state
		switch st := currentState.(type) {
		case *nfa.AcceptingState:
			// Found a match!
			return currentTask.thread.captures, true
		case *nfa.MatcherState:
			if currentTask.thread.lineIndex < len(line) {
				r, size := utf8.DecodeRune(line[currentTask.thread.lineIndex:])
				match, _ := st.Matcher.Match(r)
				if match {
					nextThread := thread{
						state:     st.Out,
						lineIndex: currentTask.thread.lineIndex + size,
						captures:  currentTask.thread.captures,
					}
					stack = append(stack, task{isRevert: false, thread: nextThread})
				}
			}
		case *nfa.SplitState:
			// Push both branches onto the stack. The order determines preference.
			thread1 := thread{state: st.Branch1, lineIndex: currentTask.thread.lineIndex, captures: currentTask.thread.captures}
			thread2 := thread{state: st.Branch2, lineIndex: currentTask.thread.lineIndex, captures: currentTask.thread.captures}
			stack = append(stack, task{isRevert: false, thread: thread2})
			stack = append(stack, task{isRevert: false, thread: thread1})
		case *nfa.CaptureStartState:
			groupIndex := st.GroupIndex
			captureSlot := groupIndex * 2

			// Log the change for backtracking.
			undo := undoEntry{captureIndex: captureSlot, oldValue: currentTask.thread.captures[captureSlot]}

			// Update the capture.
			currentTask.thread.captures[captureSlot] = currentTask.thread.lineIndex

			nextThread := thread{state: st.Out, lineIndex: currentTask.thread.lineIndex, captures: currentTask.thread.captures}

			// Push the REVERT task first, then the EXPLORE task.
			stack = append(stack, task{isRevert: true, thread: currentTask.thread, undoLog: []undoEntry{undo}})
			stack = append(stack, task{isRevert: false, thread: nextThread})

		case *nfa.CaptureEndState:
			groupIndex := st.GroupIndex
			captureSlot := groupIndex*2 + 1

			undo := undoEntry{captureIndex: captureSlot, oldValue: currentTask.thread.captures[captureSlot]}
			currentTask.thread.captures[captureSlot] = currentTask.thread.lineIndex

			nextThread := thread{state: st.Out, lineIndex: currentTask.thread.lineIndex, captures: currentTask.thread.captures}

			stack = append(stack, task{isRevert: true, thread: currentTask.thread, undoLog: []undoEntry{undo}})
			stack = append(stack, task{isRevert: false, thread: nextThread})

		case *nfa.StartAnchorState:
			if currentTask.thread.lineIndex == 0 {
				nextThread := thread{state: st.Out, lineIndex: currentTask.thread.lineIndex, captures: currentTask.thread.captures}
				stack = append(stack, task{isRevert: false, thread: nextThread})
			}
		case *nfa.EndAnchorState:
			if currentTask.thread.lineIndex == len(line) {
				nextThread := thread{state: st.Out, lineIndex: currentTask.thread.lineIndex, captures: currentTask.thread.captures}
				stack = append(stack, task{isRevert: false, thread: nextThread})
			}
		}
	}

	return nil, false
}
