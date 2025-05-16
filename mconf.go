package mconf

import (
	"io"
	"os"
	"path/filepath"

	"github.com/marzeq/mconf/parser"
	"github.com/marzeq/mconf/tokeniser"
)

func ParseFromString(s string, rootDir string, rootFile string, relativeDir string) (map[string]parser.ParserValue, map[string]parser.ParserValue, error) {
	t := tokeniser.NewTokeniser(s, rootFile, relativeDir)
	tokens, err := t.Tokenise()
	if err != nil {
		return nil, nil, err
	}

	p := parser.NewParser(tokens, rootDir, rootFile, relativeDir)
	parsed, err := p.Parse()
	if err != nil {
		return nil, nil, err
	}
	constants := p.GetConstants()

	return parsed, constants, nil
}

func ParseFromFile(filename string) (map[string]parser.ParserValue, map[string]parser.ParserValue, error) {
	f, err := os.ReadFile(filename)
	if err != nil {
		return nil, nil, err
	}

	s := string(f)

	relativeDir := filepath.Dir(filename)

	fileDir, fdirErr := filepath.Abs(relativeDir)
	if fdirErr != nil {
		return nil, nil, fdirErr
	}

	baseFile := filepath.Base(filename)

	return ParseFromString(s, fileDir, baseFile, relativeDir)
}

func ParseFromStdin() (map[string]parser.ParserValue, map[string]parser.ParserValue, error) {
	b, err := readStdin()
	if err != nil {
		return nil, nil, err
	}

	s := string(b)

	cwd, cwdErr := os.Getwd()
	if cwdErr != nil {
		return nil, nil, cwdErr
	}

	absCwd, absCwdErr := filepath.Abs(cwd)
	if absCwdErr != nil {
		return nil, nil, absCwdErr
	}

	return ParseFromString(s, absCwd, "", cwd)
}

func readStdin() ([]byte, error) {
	b, err := io.ReadAll(os.Stdin)
	if err != nil {
		return nil, err
	}

	return b, nil
}
