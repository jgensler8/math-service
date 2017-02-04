package shared

import (
  "strings"
  "strconv"
  "net/http"
  "encoding/json"
  "github.com/golang/glog"
)

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

// Equation represents a list of arguments and operators
type Equation struct {
  Arguments []Argument `json:"arguements"`
  Operators []OperatorKey `json:"operators"`
}

func handleCompute(w http.ResponseWriter, r *http.Request) {
  glog.V(2).Infof("Handling Compute Request")
  s := r.URL.Query().Get(argumentsKey)
  if len(s) == 0 {
    glog.Errorf("Bad Requests")
    w.WriteHeader(http.StatusBadRequest)
    return
  }
}

func stringToArgument(s string) (Argument) {
  i, _ := strconv.Atoi(s)
  return Argument{Value: i}
}

const argumentsKey = "args"
// GenerateHandler is a generic way to generate a handler for an Operator
func GenerateHandler(o Operator) (func(http.ResponseWriter, *http.Request)) {
  return func(w http.ResponseWriter, r *http.Request){
    glog.V(2).Infof("Handling Compute Request")
    s := r.URL.Query().Get(argumentsKey)
    if len(s) == 0 {
      glog.Errorf("Bad Requests")
      w.WriteHeader(http.StatusBadRequest)
      return
    }

    a := strings.Split(s, ",")
    if len(a) < 2 {
      glog.Errorf("Bad Requests")
      w.WriteHeader(http.StatusBadRequest)
      return
    }

    res := o.Operate(stringToArgument(a[0]), stringToArgument(a[1]))

    b, err := json.Marshal(res)
    if err != nil {
      glog.Error("Could not Marshal")
      w.WriteHeader(http.StatusInternalServerError)
      return
    }
    w.Write(b)
  }
}
