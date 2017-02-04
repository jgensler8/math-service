package main

import (
  "github.com/jgensler8/math-service/shared"
  service "github.com/jgensler8/math-service/generic-service"
)

type subtractionOperator struct {
  shared.Operator
}

func (a subtractionOperator) Operate( x shared.Argument, y shared.Argument) (shared.Argument) {
  return shared.Argument{Value: (x.Value - y.Value)}
}

// SubtractionOperator is a global Operator for subtraction
var SubtractionOperator = subtractionOperator{}

func main() {
  service.GenericServiceBuilder.
    AddHandler("/operate", shared.GenerateHandler(SubtractionOperator)).
    Build().
    Serve()
}
