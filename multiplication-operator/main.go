package main

import (
  "github.com/jgensler8/math-service/shared"
  service "github.com/jgensler8/math-service/generic-service"
)

type multiplicationOperator struct {
  shared.Operator
}

func (a multiplicationOperator) Operate( x shared.Argument, y shared.Argument) (shared.Argument) {
  return shared.Argument{Value: (x.Value * y.Value)}
}

// MultiplicationOperator is a global Operator for multiplication
var MultiplicationOperator = multiplicationOperator{}

func main() {
  service.GenericServiceBuilder.
    AddHandler("/operate", shared.GenerateHandler(MultiplicationOperator)).
    Build().
    Serve()
}
