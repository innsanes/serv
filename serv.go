package serv

type Serv struct {
	beforeServes []func() error
	afterServes  []func() error
	beforeStops  []func() error
	stopChannel  chan struct{}
	serveError   error
	log          Logger
}

func New(bfs ...BuildFunc) *Serv {
	ret := &Serv{
		beforeServes: make([]func() error, 0),
		afterServes:  make([]func() error, 0),
		beforeStops:  make([]func() error, 0),
		stopChannel:  make(chan struct{}),
		log:          &log{},
	}
	for _, bf := range bfs {
		bf(ret)
	}
	return ret
}

type BuildFunc func(x *Serv)

func WithLogger(l Logger) BuildFunc {
	return func(x *Serv) {
		x.log = l
	}
}

func (i *Serv) Serve(servers ...Server) {
	defer func() {
		if r := recover(); r != nil {
			i.log.Errorf("panic: %v", r)
			return
		}
	}()
	defer i.beforeStop(servers...)
	i.beforeServe(servers...)
	i.serve(servers...)
	i.afterServe(servers...)
	<-i.stopChannel
}

func (i *Serv) beforeServe(servers ...Server) {
	// 通过函数调用的方式注册的函数
	for _, f := range i.beforeServes {
		i.serveError = f()
		if i.serveError != nil {
			return
		}
	}
	// 每个注册的服务
	for _, s := range servers {
		i.serveError = s.BeforeServe()
		if i.serveError != nil {
			return
		}
	}
}

func (i *Serv) serve(servers ...Server) {
	for _, s := range servers {
		// 服务的启动需要异步
		i.serveError = s.Serve()
		if i.serveError != nil {
			return
		}
	}
}

func (i *Serv) afterServe(servers ...Server) {
	// 通过函数调用的方式注册的函数
	for _, f := range i.afterServes {
		i.serveError = f()
		if i.serveError != nil {
			return
		}
	}
	// 每个注册的服务
	for _, s := range servers {
		i.serveError = s.AfterServe()
		if i.serveError != nil {
			return
		}
	}
}

func (i *Serv) beforeStop(servers ...Server) {
	defer func() {
		// 错误处理
		if i.serveError != nil {
			i.log.Errorf("serve error: %s", i.serveError)
		}
		close(i.stopChannel)
	}()
	// 通过函数调用的方式注册的函数
	for _, f := range i.beforeStops {
		i.serveError = f()
		if i.serveError != nil {
			return
		}
	}
	// 每个注册的服务
	for _, s := range servers {
		i.serveError = s.BeforeStop()
		if i.serveError != nil {
			return
		}
	}
}

func (i *Serv) RegisterBeforeServe(f func() error) {
	i.beforeServes = append(i.beforeServes, f)
}

func (i *Serv) RegisterAfterServe(f func() error) {
	i.afterServes = append(i.afterServes, f)
}

func (i *Serv) RegisterBeforeStop(f func() error) {
	i.beforeStops = append(i.beforeStops, f)
}

func (i *Serv) ForceStop() {
	i.stopChannel <- struct{}{}
}
