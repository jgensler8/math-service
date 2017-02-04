package helper

import (
  "strings"
  "strconv"
  "unicode"
  "github.com/jgensler8/math-service/shared"
  ak "github.com/jgensler8/math-service/addition-operator/key"
  sk "github.com/jgensler8/math-service/subtraction-operator/key"
  mk "github.com/jgensler8/math-service/multiplication-operator/key"
  dk "github.com/jgensler8/math-service/division-operator/key"
)

func splitOperators (r rune) bool {
  return r == ak.AdditionOperatorKey.Value ||
    r == sk.SubtractionOperatorKey.Value ||
    r == mk.MultiplicationOperatorKey.Value ||
    r == dk.DivisionOperatorKey.Value
}

func splitDigits (r rune) bool {
  return unicode.IsDigit(r)
}

func extractOperators(s string) (o []shared.OperatorKey) {
  for _, t := range strings.FieldsFunc(s, splitDigits) {
    if len(t) != 0 {
      o = append(o, shared.OperatorKey{Value: rune(t[0])})
    }
  }
  return
}

func extractArguements(s string) (a []shared.Argument) {
  for _, t := range strings.FieldsFunc(s, splitOperators) {
    if len(t) != 0 {
      i, err := strconv.Atoi(t)
      if err == nil {
        a = append(a, shared.Argument{Value: i})
      }
    }
  }
  return
}

// Tokenize turns an equation string into an Equation struct
func Tokenize(s string) (e shared.Equation, err error) {
  e.Arguments = extractArguements(s)
  e.Operators = extractOperators(s)
  return
}
