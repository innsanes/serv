package serv

import "fmt"

type log struct {
}

func (l *log) Errorf(format string, v ...interface{}) {
	fmt.Println(fmt.Sprintf("ERR: "+format, v...))
}
