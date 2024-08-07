package serv

type Server interface {
	BeforeServe() error
	AfterServe() error
	Serve() error
	BeforeStop() error
}

type Logger interface {
	Errorf(format string, v ...interface{})
}
