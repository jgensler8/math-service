package main

import (
  "net/http"
  "encoding/json"
  t "github.com/jgensler8/math-service/tokenizer/helper"
  service "github.com/jgensler8/math-service/generic-service"
  "github.com/golang/glog"
)

const tokenizeKey = "string"

func handleTokenize(w http.ResponseWriter, r *http.Request) {
  glog.V(2).Infof("Handling Compute Request")
  s := r.URL.Query().Get(tokenizeKey)
  if len(s) == 0 {
    glog.Errorf("Bad Request")
    w.WriteHeader(http.StatusBadRequest)
    return
  }
  glog.V(4).Infof("Tokenizing %v", s)
  e, err := t.Tokenize(s)
  if err != nil {
    glog.Error("Could not Tokenize")
    w.WriteHeader(http.StatusInternalServerError)
    return
  }
  b, err := json.Marshal(e)
  if err != nil {
    glog.Error("Could not Marshal")
    w.WriteHeader(http.StatusInternalServerError)
    return
  }
	w.Write(b)
}

func main() {
  service.GenericServiceBuilder.
    AddHandler("/tokenize", handleTokenize).
    Build().
    Serve()
}
