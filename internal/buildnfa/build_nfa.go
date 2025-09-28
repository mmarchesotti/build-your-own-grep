package buildnfa

import (
	"github.com/mmarchesotti/build-your-own-grep/internal/ast"
	"github.com/mmarchesotti/build-your-own-grep/internal/matcher"
	"github.com/mmarchesotti/build-your-own-grep/internal/nfa"
	"github.com/mmarchesotti/build-your-own-grep/internal/parser"
	"github.com/mmarchesotti/build-your-own-grep/internal/predefinedclass"
)

func processNode(n ast.ASTNode) nfa.Fragment {
	switch node := n.(type) {
	case *ast.AlternationNode:
		subfragment1 := processNode(node.Left)
		subfragment2 := processNode(node.Right)
		return nfa.Fragment{
			Start: &nfa.SplitState{
				Branch1: subfragment1.Start,
				Branch2: subfragment2.Start,
			},
			Out: append(subfragment1.Out, subfragment2.Out...),
		}
	case *ast.ConcatenationNode:
		subfragment1 := processNode(node.Left)
		subfragment2 := processNode(node.Right)
		nfa.SetStates(subfragment1.Out, subfragment2.Start)
		return nfa.Fragment{
			Start: subfragment1.Start,
			Out:   subfragment2.Out,
		}
	case *ast.KleeneClosureNode:
		subfragment := processNode(node.Child)
		split := nfa.SplitState{
			Branch1: subfragment.Start,
			Branch2: nil,
		}
		nfa.SetStates(subfragment.Out, &split)
		return nfa.Fragment{
			Start: &split,
			Out:   []*nfa.State{&split.Branch2},
		}
	case *ast.PositiveClosureNode:
		subfragment := processNode(node.Child)
		split := nfa.SplitState{
			Branch1: subfragment.Start,
			Branch2: nil,
		}
		nfa.SetStates(subfragment.Out, &split)
		return nfa.Fragment{
			Start: subfragment.Start,
			Out:   []*nfa.State{&split.Branch2},
		}
	case *ast.OptionalNode:
		subfragment := processNode(node.Child)
		split := nfa.SplitState{
			Branch1: subfragment.Start,
			Branch2: nil,
		}
		return nfa.Fragment{
			Start: &split,
			Out:   append(subfragment.Out, &split.Branch2),
		}
	case *ast.LiteralNode:
		matcher := nfa.MatcherState{
			Out:     nil,
			Matcher: &matcher.LiteralMatcher{Literal: node.Literal},
		}
		return nfa.Fragment{
			Start: &matcher,
			Out:   []*nfa.State{&matcher.Out},
		}
	case *ast.CharacterSetNode:
		var characterClassesMatchers []matcher.PredefinedClassMatcher
		for _, characterClass := range node.CharacterClasses {
			var m matcher.PredefinedClassMatcher
			switch characterClass {
			case predefinedclass.ClassDigit:
				m = &matcher.DigitMatcher{}
			case predefinedclass.ClassAlphanumeric:
				m = &matcher.AlphaNumericMatcher{}
			}
			characterClassesMatchers = append(characterClassesMatchers, m)
		}
		matcher := nfa.MatcherState{
			Out: nil,
			Matcher: &matcher.CharacterSetMatcher{
				IsPositive:               node.IsPositive,
				Literals:                 node.Literals,
				Ranges:                   node.Ranges,
				CharacterClassesMatchers: characterClassesMatchers,
			},
		}
		return nfa.Fragment{
			Start: &matcher,
			Out:   []*nfa.State{&matcher.Out},
		}
	case *ast.WildcardNode:
		matcher := nfa.MatcherState{
			Out:     nil,
			Matcher: &matcher.WildcardMatcher{},
		}
		return nfa.Fragment{
			Start: &matcher,
			Out:   []*nfa.State{&matcher.Out},
		}
	case *ast.DigitNode:
		matcher := nfa.MatcherState{
			Out:     nil,
			Matcher: &matcher.DigitMatcher{},
		}
		return nfa.Fragment{
			Start: &matcher,
			Out:   []*nfa.State{&matcher.Out},
		}
	case *ast.AlphaNumericNode:
		matcher := nfa.MatcherState{
			Out:     nil,
			Matcher: &matcher.AlphaNumericMatcher{},
		}
		return nfa.Fragment{
			Start: &matcher,
			Out:   []*nfa.State{&matcher.Out},
		}
	case *ast.StartAnchorNode:
		s := &nfa.StartAnchorState{
			Out: nil,
		}
		return nfa.Fragment{
			Start: s,
			Out:   []*nfa.State{&s.Out},
		}
	case *ast.EndAnchorNode:
		s := &nfa.EndAnchorState{
			Out: nil,
		}
		return nfa.Fragment{
			Start: s,
			Out:   []*nfa.State{&s.Out},
		}
	default:
		return nfa.Fragment{}
		// TODO ERROR
	}
}

func Build(inputPattern string) nfa.Fragment {
	tree := parser.Parse(inputPattern)
	f := processNode(tree)
	acceptingState := &nfa.AcceptingState{}
	nfa.SetStates(f.Out, acceptingState)
	f.Out = []*nfa.State{}
	return f
}
