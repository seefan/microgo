# microgo
The Microservice Framework develop by go

### example

    import (
        "github.com/seefan/microgo/server"
        "github.com/seefan/microgo/server/thriftworker"
        "git.apache.org/thrift.git/lib/go/thrift"
        "github.com/seefan/microgo/test/gen-go/test"
        test2 "github.com/seefan/microgo/test"
    )
    
    func main() {
        //define a tcp worker
        run := thriftworker.NewTcpWorker()
        //register all thrift processor
        run.RegisterThriftProcessor("test.HelloWorld", func() thrift.TProcessor {
        	return test.NewHelloWorldProcessor(&test2.HelloWorldImpl{})
        })
        //define transport and protocol,default is framed,binary
        //run.TransportFactory = thrift.NewTFramedTransportFactory(thrift.NewTTransportFactory())
        //run.ProtocolFactory = thrift.NewTBinaryProtocolFactoryDefault()
        //Register the correspondence between service name and service id to reduce network traffic
        server.RegisterServiceId("test.HelloWorld", "1002")
        //run the worker
        server.Run(run)
    }
