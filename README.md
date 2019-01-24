# microgo
The Microservice Framework develop by go

### example

#### main

    package main

    import (
        "context"
        "github.com/seefan/microgo/server/httpserver"
        "github.com/seefan/microgo/test"
    )
    
    func main() {
        s := httpserver.NewHTTPServer("localhost", 8888)
        s.Register(&test.TestService{})
        if err := s.Start(context.Background()); err != nil {
            println(err.Error())
        }
    
    }


#### service

    package test
    
    import (
        "github.com/seefan/microgo/service"
    )
    
    type TestService struct {
    }
    
    func (TestService) Hello(entry service.Entry) string {
        name := entry.Get("name")
        return "hello " + name
    }
    func (TestService) Name() string {
        return "test"
    }
    func (TestService) Version() string {
        return "1.0"
    }
