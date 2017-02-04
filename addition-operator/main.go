package main

import (
  "github.com/jgensler8/math-service/shared"
  service "github.com/jgensler8/math-service/generic-service"
)

type additionOperator struct {
  shared.Operator
}

func (a additionOperator) Operate( x shared.Argument, y shared.Argument) (shared.Argument) {
  return shared.Argument{Value: (x.Value + y.Value)}
}

// AdditionOperator is a global Operator for addition
var AdditionOperator = additionOperator{}

func main() {
  service.GenericServiceBuilder.
    AddHandler("/operate", shared.GenerateHandler(AdditionOperator)).
    Build().
    Serve()
}
