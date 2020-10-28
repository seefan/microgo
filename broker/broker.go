package broker

//Broker service broker support
type Broker interface {
	//Register register service
	//name servie name
	//host service host
	//port service port
	Register(name, path, host string, port int) error
}
