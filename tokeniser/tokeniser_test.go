package tokeniser_test

import (
	"testing"

	"github.com/marzeq/mconf/tokeniser"
)

func checkTokenEqual(t *testing.T, token tokeniser.Token, expectedToken tokeniser.Token) {
  if token.Type != expectedToken.Type {
    t.Fatalf("Expected token type %s, got %s", expectedToken.Type, token.Type)
  }

  if token.Value != expectedToken.Value {
    t.Fatalf("Expected token value %s, got %s", expectedToken.Value, token.Value)
  }

  if token.Start.Line != expectedToken.Start.Line {
    t.Fatalf("Expected token start line %d, got %d", expectedToken.Start.Line, token.Start.Line)
  }

  if token.Start.Col != expectedToken.Start.Col {
    t.Fatalf("Expected token start col %d, got %d", expectedToken.Start.Col, token.Start.Col)
  }

  if token.Start.Line != expectedToken.Start.Line {
    t.Fatalf("Expected token start line %d, got %d", expectedToken.Start.Line, token.Start.Line)
  }
}

func checkTokensEqual(t *testing.T, tokens []tokeniser.Token, expectedTokens []tokeniser.Token) {
  if len(tokens) != len(expectedTokens) {
    t.Fatalf("Expected %d tokens, got %d", len(expectedTokens), len(tokens))
  }

  for i, token := range tokens {
    checkTokenEqual(t, token, expectedTokens[i])
  }
}

func TestEmpty(t *testing.T) {
  tks := tokeniser.NewTokeniser("")

  tokens := tks.Tokenise()

  checkTokensEqual(t, tokens, []tokeniser.Token{})
}

func TestNumberValue(t *testing.T) {
  tks := tokeniser.NewTokeniser("123")

  tokens := tks.Tokenise()

  checkTokensEqual(t, tokens, []tokeniser.Token{
    tokeniser.NumberToken("123", tokeniser.Location{Line: 1, Col: 1}),
  })
}

func TestStringValue(t *testing.T) {
  tks := tokeniser.NewTokeniser("\"abc abc \\\"abc\\\" 'abc'\n123\"")
  
  tokens := tks.Tokenise()

  checkTokensEqual(t, tokens, []tokeniser.Token{
    tokeniser.StringToken("abc abc \"abc\" 'abc'\n123", tokeniser.Location{Line: 1, Col: 1}),
  })
}

func TestAssignment(t *testing.T) {
  tks := tokeniser.NewTokeniser("abc = 123")

  tokens := tks.Tokenise()

  checkTokensEqual(t, tokens, []tokeniser.Token{
    tokeniser.KeyToken("abc", tokeniser.Location{Line: 1, Col: 1}),
    tokeniser.AssignToken(tokeniser.Location{Line: 1, Col: 5}),
    tokeniser.NumberToken("123", tokeniser.Location{Line: 1, Col: 7}),
  })
}

func TestComment(t *testing.T) {
  tks := tokeniser.NewTokeniser("abc # 123")

  tokens := tks.Tokenise()

  checkTokensEqual(t, tokens, []tokeniser.Token{
    tokeniser.KeyToken("abc", tokeniser.Location{Line: 1, Col: 1}),
  })
}

func TestConstantAssignment(t *testing.T) {
  tks := tokeniser.NewTokeniser("$val = 123")

  tokens := tks.Tokenise()

  checkTokensEqual(t, tokens, []tokeniser.Token{
    tokeniser.ConstantToken("val", tokeniser.Location{Line: 1, Col: 1}),
    tokeniser.AssignToken(tokeniser.Location{Line: 1, Col: 6}),
    tokeniser.NumberToken("123", tokeniser.Location{Line: 1, Col: 8}),
  })
}

func TestList(t *testing.T) {
  tks := tokeniser.NewTokeniser("[1, 2]")

  tokens := tks.Tokenise()

  checkTokensEqual(t, tokens, []tokeniser.Token{
    tokeniser.OpenListToken(tokeniser.Location{Line: 1, Col: 1}),
    tokeniser.NumberToken("1", tokeniser.Location{Line: 1, Col: 2}),
    tokeniser.CommaToken(tokeniser.Location{Line: 1, Col: 3}),
    tokeniser.NumberToken("2", tokeniser.Location{Line: 1, Col: 5}),
    tokeniser.CloseListToken(tokeniser.Location{Line: 1, Col: 6}),
  })
}

func TestObject(t *testing.T) {
  tks := tokeniser.NewTokeniser("{\n  key1 = 1\n  key2 = 2\n}")

  tokens := tks.Tokenise()

  checkTokensEqual(t, tokens, []tokeniser.Token{
    tokeniser.OpenObjToken(tokeniser.Location{Line: 1, Col: 1}),
    tokeniser.KeyToken("key1", tokeniser.Location{Line: 2, Col: 3}),
    tokeniser.AssignToken(tokeniser.Location{Line: 2, Col: 8}),
    tokeniser.NumberToken("1", tokeniser.Location{Line: 2, Col: 10}),
    tokeniser.KeyToken("key2", tokeniser.Location{Line: 3, Col: 3}),
    tokeniser.AssignToken(tokeniser.Location{Line: 3, Col: 8}),
    tokeniser.NumberToken("2", tokeniser.Location{Line: 3, Col: 10}),
    tokeniser.CloseObjToken(tokeniser.Location{Line: 4, Col: 1}),
  })
}
