package helper

import (
  "testing"
  "github.com/jgensler8/math-service/shared"
)

var validTestString = "1+2*3/4234-9"

func TestTokenize(t *testing.T) {
  var e shared.Equation
  var err error
  e, err = Tokenize(validTestString)
  if err != nil {
    t.Errorf("Expected no error but got %s", err.Error())
  }
  if( len(e.Operators) != 4) {
    t.Errorf("Expected %d Operators but got %d. %v", 4, len(e.Operators), e)
  }
  if( len(e.Arguments) != 5) {
    t.Errorf("Expected %d Arguments but got %d. %v", 5, len(e.Arguments), e)
  }
}
