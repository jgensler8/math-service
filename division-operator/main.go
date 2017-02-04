package main

import (
  "github.com/jgensler8/math-service/shared"
  service "github.com/jgensler8/math-service/generic-service"
)

type divisionOperator struct {
  shared.Operator
}

func (a divisionOperator) Operate( x shared.Argument, y shared.Argument) (shared.Argument) {
  return shared.Argument{Value: (x.Value / y.Value)}
}

// DivisionOperator is a global Operator for division
var DivisionOperator = divisionOperator{}

func main() {
  service.GenericServiceBuilder.
    AddHandler("/operate", shared.GenerateHandler(DivisionOperator)).
    Build().
    Serve()
}
