package serv

var serv = New()

func Serve(servers ...Server) {
	serv.Serve(servers...)
}

func RegisterBeforeServe(f func() error) {
	serv.beforeServes = append(serv.beforeServes, f)
}

func RegisterAfterServe(f func() error) {
	serv.afterServes = append(serv.afterServes, f)
}

func RegisterBeforeStop(f func() error) {
	serv.beforeStops = append(serv.beforeStops, f)
}

func ForceStop() {
	serv.ForceStop()
}
