package backtrack

// import (
// 	"github.com/mmarchesotti/build-your-own-grep/internal/nfasimulator"
// 	"github.com/mmarchesotti/build-your-own-grep/internal/token"
// )

// func Run(line []byte, tokens []token.Token) (match bool, err error) {
// 	return processTokens(line, 0, tokens, []nfasimulator.Capture{})
// }

// func processTokens(line []byte, lineIndex int, tokens []token.Token, allCapturedGroups []nfasimulator.Capture) (match bool, err error) {
// 	for i, t := range tokens {
// 		backReferenceToken, isBackreferenceToken := t.(*token.BackReference)
// 		if !isBackreferenceToken {
// 			continue
// 		}

// 		matchesChannel := nfasimulator.Simulate(line[lineIndex:], tokens[:i])
// 		for match := range matchesChannel {
// 			localLineIndex := lineIndex
// 			allCapturedGroups = append(allCapturedGroups, match.capturedGroups)
// 			backReferenceGroup, err := getReferencedGroup(allCapturedGroups, backReferenceToken)
// 			if err != nil {
// 				return false, err
// 			}
// 			backReferenceMatch, err := matchBackReference(line, localLineIndex, backReferenceGroup)
// 			if err != nil {
// 				return false, err
// 			}
// 			if !backReferenceMatch {
// 				continue
// 			}
// 			localLineIndex += len(backReferenceGroup)
// 			restOfPatternMatch, err := processTokens(line, localLineIndex, tokens[i+1:], captureGroup)
// 			if err != nil {
// 				return false, err
// 			}
// 			if restOfPatternMatch {
// 				return true, nil
// 			}
// 		}
// 		return false, nil
// 	}
// 	return Simulate(line[lineIndex:], tokens)
// }

// func getReferencedGroup(allCapturedGroups [][]byte, backReferenceToken BackReferenceToken) ([]byte, error) {
// 	if backReferenceToken.CapturedGroupNumber > len(allCapturedGroups) {
// 		return nil, fmt.Error("Reference to non-existing group '%d'", backReferenceToken.CapturedGroupNumber)
// 	}
// 	capturedGroupIndex := backReferenceToken.CapturedGroupNumber - 1
// 	return allCapturedGroups[capturedGroupIndex], nil
// }
