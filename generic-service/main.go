package genericservice

import (
  "fmt"
  "net/http"
  "flag"
  "github.com/golang/glog"
  "github.com/kelseyhightower/envconfig"
  "github.com/lann/builder"
)

// GenericService is a strcut that holds configuration on how to build a common
// Dockerize-able go http service
type GenericService struct {
  Specification interface{}
}

type genericServiceBuilder builder.Builder

func (b genericServiceBuilder) AddHandler(e string, h func(http.ResponseWriter, *http.Request)) genericServiceBuilder {
  http.HandleFunc(e, h)
  return b
}

func (b genericServiceBuilder) AddSpecification(s interface{}) genericServiceBuilder {
    return builder.Set(b, "Specification", s).(genericServiceBuilder)
}

func (b genericServiceBuilder) Build() GenericService {
    return builder.GetStruct(b).(GenericService)
}

// GenericServiceBuilder is a builder for services
var GenericServiceBuilder = builder.Register(genericServiceBuilder{}, GenericService{}).(genericServiceBuilder)

var port int

func (g GenericService) init() {
  flag.IntVar(&port, "port", 8080, "port to listen on")
  flag.Parse()
  if g.Specification != nil {
    err := envconfig.Process("gateway", g.Specification)
    if err != nil {
        glog.Fatal(err.Error())
    }
  }
  glog.V(1).Infof("Listening on port %d", port)
}

// Serve runs the service as an http service
func (g GenericService) Serve() {
  g.init()
  glog.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}
