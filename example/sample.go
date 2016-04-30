package main

import (
	"fmt"
	"github.com/gophil/npd"
	_ "net/http/pprof"
	"runtime"
)

type MyTask struct {
	number  int
	message string
}

func NewMyTask(number int, message string) *MyTask {
	return &MyTask{
		number:  number,
		message: message,
	}
}

func (m *MyTask) DoSNMP() error {
	fmt.Println(m.message, " -------> ", m.number)
	return nil
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	d := npd.NewDispatcher(30, 30)

	d.Run()
	defer d.Stop()

	for i := 0; i < 30; i++ {
		targetObject := NewMyTask(i, "execute demo")
		task := npd.CreateTask(targetObject, "DoSNMP")
		d.SubmitTask(task)
	}

	select {}

}
