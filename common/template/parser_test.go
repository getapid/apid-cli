package template

import (
	"testing"
)

type parseTest struct {
	input  string
	tokens []token
}

func mkToken(typ tokenType, text string) token {
	return token{
		typ: typ,
		val: text,
	}
}

var (
	tEnd = mkToken(tokenEnd, "")
)

var parseTests = []parseTest{
	{
		"",
		[]token{
			tEnd,
		},
	},
	{
		"some text",
		[]token{
			mkToken(tokenText, "some text"),
			tEnd,
		},
	},
	{
		"{{ key }}",
		[]token{
			mkToken(tokenIdentifier, "key"),
			tEnd,
		},
	},
	{
		"{% command %}",
		[]token{
			mkToken(tokenCommand, "command"),
			tEnd,
		},
	},
	{
		"pre {{ key }}",
		[]token{
			mkToken(tokenText, "pre "),
			mkToken(tokenIdentifier, "key"),
			tEnd,
		},
	},
	{
		"pre {% command %} post",
		[]token{
			mkToken(tokenText, "pre "),
			mkToken(tokenCommand, "command"),
			mkToken(tokenText, " post"),
			tEnd,
		},
	},
	{
		"pre {{ key }} between {% command %} post",
		[]token{
			mkToken(tokenText, "pre "),
			mkToken(tokenIdentifier, "key"),
			mkToken(tokenText, " between "),
			mkToken(tokenCommand, "command"),
			mkToken(tokenText, " post"),
			tEnd,
		},
	},
	{
		"{{ key }} post",
		[]token{
			mkToken(tokenIdentifier, "key"),
			mkToken(tokenText, " post"),
			tEnd,
		},
	},
	{
		"some {{ key }} text",
		[]token{
			mkToken(tokenText, "some "),
			mkToken(tokenIdentifier, "key"),
			mkToken(tokenText, " text"),
			tEnd,
		},
	},
	{
		"some {{ key }} text {{ other_key }}",
		[]token{
			mkToken(tokenText, "some "),
			mkToken(tokenIdentifier, "key"),
			mkToken(tokenText, " text "),
			mkToken(tokenIdentifier, "other_key"),
			tEnd,
		},
	},
}

// collect gathers the emitted items into a slice.
func collectTokens(t parseTest) (tokens []token) {
	p := parse(t.input)
	for {
		token := p.nextItem()
		tokens = append(tokens, token)
		if token.typ == tokenEnd || token.typ == tokenError {
			break
		}
	}
	return
}

func tokensEqual(i1, i2 []token) bool {
	if len(i1) != len(i2) {
		return false
	}
	for k := range i1 {
		if i1[k].typ != i2[k].typ {
			return false
		}
		if i1[k].val != i2[k].val && i1[k].typ != tokenError {
			return false
		}
	}
	return true
}

func TestParse(t *testing.T) {
	for _, test := range parseTests {
		tokens := collectTokens(test)
		if !tokensEqual(tokens, test.tokens) {
			t.Errorf("got\n\t%+v\nexpected\n\t%v", tokens, test.tokens)
		}
	}
}
