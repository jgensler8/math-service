package main

import (
  "errors"
  "strconv"
  "net/http"
  "net/url"
  "encoding/json"
  "github.com/golang/glog"
  "github.com/jgensler8/math-service/shared"
  ak "github.com/jgensler8/math-service/addition-operator/key"
  sk "github.com/jgensler8/math-service/subtraction-operator/key"
  mk "github.com/jgensler8/math-service/multiplication-operator/key"
  dk "github.com/jgensler8/math-service/division-operator/key"
  service "github.com/jgensler8/math-service/generic-service"
)

// Specification hold the services we integrate with. These should be urls
type Specification struct {
  TokenizerService      string `default:"http://localhost:8081/tokenize"`
  AdditionService       string `default:"http://localhost:8082/operate"`
  SubtractionService    string `default:"http://localhost:8083/operate"`
  MultiplicationService string `default:"http://localhost:8084/operate"`
  DivisionService       string `default:"http://localhost:8085/operate"`
}

const equationKey = "equation"

var spec Specification

func operatorKeyToService(k shared.OperatorKey) (string, error){
  switch k {
  case ak.AdditionOperatorKey:
    return spec.AdditionService, nil
  case sk.SubtractionOperatorKey:
    return spec.SubtractionService, nil
  case mk.MultiplicationOperatorKey:
    return spec.MultiplicationService, nil
  case dk.DivisionOperatorKey:
    return spec.DivisionService, nil
  }
  return "", errors.New("Unsupported Operator")
}

func serviceRequest(u url.URL, i interface{}) (err error) {
  glog.V(4).Infof("Sending Request to Service (%s)", u.String())
  resp, err := http.Get(u.String())
  if err != nil {
    glog.Errorf("Couldn't reach Service: %s", err.Error())
    return
  }

  glog.V(4).Infof("Decoding Response from Service")
  err = json.NewDecoder(resp.Body).Decode(i)
  if err != nil {
    glog.Errorf("Couldn't decode response from Service (%v)", i)
    return
  }

  return
}

func constructServiceURL(s string, k string, v string) (u *url.URL, err error) {
  glog.V(5).Infof("Building Service Request (%s) (%s) (%s)", s, k, v)
  u, err = url.Parse(s)
  if err != nil {
    glog.Errorf("%v", err.Error())
    return nil, err
  }
  q := u.Query()
  q.Set(k, v)
  u.RawQuery = q.Encode()
  return u, nil
}

func operate(k shared.OperatorKey, x shared.Argument, y shared.Argument) (z shared.Argument, err error) {
  glog.V(6).Infof("%v %v %v", x.Value, k.Value, y.Value)
  s, err := operatorKeyToService(k)
  if err != nil {
    return z, err
  }

  v := strconv.Itoa(x.Value) + "," + strconv.Itoa(y.Value)
  u, err := constructServiceURL(s, "args", v)
  if err != nil {
    return z, err
  }

  err = serviceRequest(*u, &z)
  glog.V(6).Infof("Response: %v", z.Value)
  return z, err
}

func handleCompute(w http.ResponseWriter, r *http.Request) {
  glog.V(2).Infof("Handling Compute Request")
  s := r.URL.Query().Get(equationKey)
  if len(s) == 0 {
    glog.Errorf("Bad Requests")
    w.WriteHeader(http.StatusBadRequest)
    return
  }

  u, err := constructServiceURL(spec.TokenizerService, "string", s)
  if err != nil {
    glog.Errorf("%v", err.Error())
    w.WriteHeader(http.StatusInternalServerError)
    return
  }

  e := shared.Equation{}
  glog.V(2).Infof("Sending request to Tokenizer Service")
  err = serviceRequest(*u, &e)
  if err != nil {
    glog.Errorf("%v", err.Error())
    w.WriteHeader(http.StatusInternalServerError)
    return
  }

  var last = e.Arguments[0]
  for i := 1; i < len(e.Arguments); i++ {
    last, err = operate(e.Operators[i-1], last, e.Arguments[i])
    if err != nil {
      w.WriteHeader(http.StatusInternalServerError)
      return
    }
  }
  glog.V(4).Infof("%d", last.Value)
  w.Write([]byte(strconv.Itoa(last.Value)))
}

func main() {
  service.GenericServiceBuilder.
    AddSpecification(&spec).
    AddHandler("/compute", handleCompute).
    Build().
    Serve()
}
