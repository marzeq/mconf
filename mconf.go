package mconf

import (
  "io"
  "os"
  "path/filepath"

  "github.com/marzeq/mconf/mconf_parser"
  "github.com/marzeq/mconf/mconf_tokeniser"
  "github.com/marzeq/mconf/mconf_values"
)

const (
  VERSION  = "2.0.0"
  PROGNAME = "mconf"
)

func ParseFromString(s string, rootDir string, rootFile string, relativeDir string, constsOpt ...map[string]mconf_values.MconfValue) (*mconf_values.MconfObject, map[string]mconf_values.MconfValue, error) {
  t := mconf_tokeniser.NewTokeniser(s, rootFile, relativeDir)
  tokens, err := t.Tokenise()
  if err != nil {
    return nil, nil, err
  }

  var consts map[string]mconf_values.MconfValue
  if len(constsOpt) > 0 {
    consts = constsOpt[0]
  } else {
    consts = make(map[string]mconf_values.MconfValue)
  }

  p := mconf_parser.NewParser(tokens, rootDir, rootFile, relativeDir, consts)
  parsed, err := p.Parse()
  if err != nil {
    return nil, nil, err
  }
  constants := p.Constants

  return parsed, constants, nil
}

func ParseFromFile(filename string, constsOpt ...map[string]mconf_values.MconfValue) (*mconf_values.MconfObject, map[string]mconf_values.MconfValue, error) {
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

  if len(constsOpt) > 0 {
    return ParseFromString(s, fileDir, baseFile, relativeDir, constsOpt[0])
  }

  return ParseFromString(s, fileDir, baseFile, relativeDir)
}

func ParseFromStdin(constsOpt ...map[string]mconf_values.MconfValue) (*mconf_values.MconfObject, map[string]mconf_values.MconfValue, error) {
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

  if len(constsOpt) > 0 {
    return ParseFromString(s, absCwd, "", cwd, constsOpt[0])
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
