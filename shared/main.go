package shared

import (
  "flag"
  "errors"
  "strings"
  "strconv"
  "net/http"
  "encoding/json"
  "github.com/golang/glog"
  "github.com/eawsy/aws-lambda-go-core/service/lambda/runtime"
)

func init() {
  flag.Set("v", "2")
  flag.Set("v", "2")
}

// Argument is an integer
type Argument struct {
  Value int `json:"value"`
}

// OperatorKey represents a way to turn
type OperatorKey struct {
  Value rune `json:"value"`
}

// Operator represents an object that can operate on Arugements
type Operator interface {
  Operate(Argument, Argument) Argument
}

// LambdaOperatorQueryStringParameters are query parameters specifc to Math service operators that are hosted on AWS Lambda
type LambdaOperatorQueryStringParameters struct {
  Arguments string `json:"args"`
}

// LambdaOperatorEvent is a json struct constructed by AWS Lambda when our Lambda function is invoked.
type LambdaOperatorEvent struct {
  QueryStringParameters LambdaOperatorQueryStringParameters `json:"queryStringParameters"`
}

/*
{
    "isBase64Encoded": true|false,
    "statusCode": httpStatusCode,
    "headers": { "headerName": "headerValue", ... },
    "body": "..."
}
*/
// LambdaProxyResponseFormat is the generic response that Lambda expects *when using proxy mode*. See http://docs.aws.amazon.com/apigateway/latest/developerguide/set-up-lambda-proxy-integrations.html
type LambdaProxyResponseFormat struct {
  IsBase64Encoded bool `json:"isBase64Encoded"`
  StatusCode int `json:"statusCode"`
  Headers map[string]string `json:"headers,omitempty"`
  Body string `json:"body"`
}

// Equation represents a list of arguments and operators
type Equation struct {
  Arguments []Argument `json:"arguements"`
  Operators []OperatorKey `json:"operators"`
}

func stringToArgument(s string) (Argument) {
  i, _ := strconv.Atoi(s)
  return Argument{Value: i}
}

const argumentsKey = "args"

// QueryStringToArguments will convert a query string to an Operator into two Arguments
func QueryStringToArguments(s string) (arg1 Argument, arg2 Argument, err error) {
  if len(s) == 0 {
    glog.Errorf("Bad Request")
    return arg1, arg2, errors.New("Bad Request: zero arguments")
  }

  a := strings.Split(s, ",")
  if len(a) < 2 {
    glog.Errorf("Bad Requests")
    return arg1, arg2, errors.New("Bad Request: fewer than 2 arguments")
  }
  return stringToArgument(a[0]), stringToArgument(a[1]), nil
}

// GenerateHandler is a generic way to generate a handler for an Operator
func GenerateHandler(o Operator) (func(http.ResponseWriter, *http.Request)) {
  return func(w http.ResponseWriter, r *http.Request){
    glog.V(2).Infof("Handling Compute Request")
    s := r.URL.Query().Get(argumentsKey)

    arg1, arg2, err := QueryStringToArguments(s)
    if err != nil {
      w.WriteHeader(http.StatusBadRequest)
      return
    }

    res := o.Operate(arg1, arg2)

    b, err := json.Marshal(res)
    if err != nil {
      glog.Error("Could not Marshal result of Operator")
      w.WriteHeader(http.StatusInternalServerError)
      return
    }
    w.Write(b)
  }
}

// GenerateEawsyLambdaHandler is a way to generate a function for the eawsy runtime in Lambda
func GenerateEawsyLambdaHandler(o Operator) ( func(evt json.RawMessage, ctx *runtime.Context)(LambdaProxyResponseFormat, error) ) {
  return func(evt json.RawMessage, ctx *runtime.Context)(res LambdaProxyResponseFormat, err error) {
    var params LambdaOperatorEvent
    err = json.Unmarshal(evt, &params)
    if err != nil {
      glog.Error("Failed to Unmarshal Operator parameters")
      return res, err
    }

    arg1, arg2, err := QueryStringToArguments(params.QueryStringParameters.Arguments)
    if err != nil {
      glog.Error("Failed to Parse Query String Arguments")
      return res, err
    }

    resultArg := o.Operate(arg1, arg2)

    res.IsBase64Encoded = false
    res.StatusCode = 200
    resultByte, _ := json.Marshal(resultArg)
    res.Body = string(resultByte)
    return res, nil
  }
}
