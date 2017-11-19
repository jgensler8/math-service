package main

import (
  "encoding/json"
  "github.com/jgensler8/math-service/shared"
  "github.com/eawsy/aws-lambda-go-core/service/lambda/runtime"
)

type additionOperator struct {
  shared.Operator
}

func (a additionOperator) Operate( x shared.Argument, y shared.Argument) (shared.Argument) {
  return shared.Argument{Value: (x.Value + y.Value)}
}

// AdditionOperator is a global Operator for addition
var AdditionOperator = additionOperator{}
var AdditionLambda func(evt json.RawMessage, ctx *runtime.Context)(shared.LambdaProxyResponseFormat, error)

func init() {
  AdditionLambda = shared.GenerateEawsyLambdaHandler(AdditionOperator)
}

func Handle(evt json.RawMessage, ctx *runtime.Context) (shared.LambdaProxyResponseFormat, error) {
  return AdditionLambda(evt, ctx)
}
