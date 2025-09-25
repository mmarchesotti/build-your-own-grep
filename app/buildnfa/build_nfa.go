package buildnfa

import (
	"github.com/mmarchesotti/build-your-own-grep/app/ast"
	"github.com/mmarchesotti/build-your-own-grep/app/nfa/frag"
	"github.com/mmarchesotti/build-your-own-grep/app/nfa/matcher"
	"github.com/mmarchesotti/build-your-own-grep/app/parser"
	"github.com/mmarchesotti/build-your-own-grep/app/predefinedclass"
)

func processNode(n ast.ASTNode) frag.Fragment {
	switch node := n.(type) {
	case *ast.AlternationNode:
		subfragment1 := processNode(node.Left)
		subfragment2 := processNode(node.Right)
		return frag.Fragment{
			Start: &frag.SplitState{
				Branch1: &subfragment1.Start,
				Branch2: &subfragment2.Start,
			},
			Out: append(subfragment1.Out, subfragment2.Out...),
		}
	case *ast.ConcatenationNode:
		subfragment1 := processNode(node.Left)
		subfragment2 := processNode(node.Right)
		frag.SetStates(subfragment1.Out, subfragment2.Start)
		return frag.Fragment{
			Start: subfragment1.Start,
			Out:   subfragment2.Out,
		}
	case *ast.KleeneClosureNode:
		subfragment := processNode(node.Child)
		split := frag.SplitState{
			Branch1: &subfragment.Start,
			Branch2: nil,
		}
		frag.SetStates(subfragment.Out, &split)
		return frag.Fragment{
			Start: &split,
			Out:   []*frag.State{split.Branch2},
		}
	case *ast.PositiveClosureNode:
		subfragment := processNode(node.Child)
		split := frag.SplitState{
			Branch1: &subfragment.Start,
			Branch2: nil,
		}
		frag.SetStates(subfragment.Out, &split)
		return frag.Fragment{
			Start: subfragment.Start,
			Out:   []*frag.State{split.Branch2},
		}
	case *ast.OptionalNode:
		subfragment := processNode(node.Child)
		split := frag.SplitState{
			Branch1: &subfragment.Start,
			Branch2: nil,
		}
		return frag.Fragment{
			Start: &split,
			Out:   append(subfragment.Out, split.Branch2),
		}
	case *ast.LiteralNode:
		matcher := frag.MatcherState{
			Out:     nil,
			Matcher: &matcher.LiteralMatcher{Literal: node.Literal},
		}
		return frag.Fragment{
			Start: &matcher,
			Out:   []*frag.State{matcher.Out},
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
		matcher := frag.MatcherState{
			Out: nil,
			Matcher: &matcher.CharacterSetMatcher{
				Literals:                 node.Literals,
				Ranges:                   node.Ranges,
				CharacterClassesMatchers: characterClassesMatchers,
			},
		}
		return frag.Fragment{
			Start: &matcher,
			Out:   []*frag.State{matcher.Out},
		}
	case *ast.WildcardNode:
		matcher := frag.MatcherState{
			Out:     nil,
			Matcher: &matcher.WildcardMatcher{},
		}
		return frag.Fragment{
			Start: &matcher,
			Out:   []*frag.State{matcher.Out},
		}
	case *ast.DigitNode:
		matcher := frag.MatcherState{
			Out:     nil,
			Matcher: &matcher.DigitMatcher{},
		}
		return frag.Fragment{
			Start: &matcher,
			Out:   []*frag.State{matcher.Out},
		}
	case *ast.AlphaNumericNode:
		matcher := frag.MatcherState{
			Out:     nil,
			Matcher: &matcher.AlphaNumericMatcher{},
		}
		return frag.Fragment{
			Start: &matcher,
			Out:   []*frag.State{matcher.Out},
		}
	case *ast.StartAnchorNode:
		return frag.Fragment{
			Start: &frag.StartAnchorState{},
		}
	case *ast.EndAnchorNode:
		return frag.Fragment{
			Start: &frag.EndAnchorState{},
		}
	default:
		return frag.Fragment{}
		// TODO ERROR
	}
}

func Build(inputPattern string) frag.Fragment {
	tree := parser.Parse(inputPattern)
	f := processNode(tree)
	acceptingState := &frag.AcceptingState{}
	frag.SetStates(f.Out, acceptingState)
	return f
}
