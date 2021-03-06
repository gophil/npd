package main

import (
	"fmt"
	"github.com/gophil/npd"
	_ "net/http/pprof"
	"runtime"
	"sync"
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

func (m *MyTask) DoSNMP() {
	fmt.Println(m.message, " -> ", m.number)
}

func main() {

	runtime.GOMAXPROCS(runtime.NumCPU())

	var wg sync.WaitGroup

	d := npd.NewDispatcherWithWait(30, 30, &wg)

	d.Run()
	defer d.Stop()

	wg.Add(1)

	go func() {
		for i := 0; i < 30; i++ {
			task := npd.CreateTask(NewMyTask(i, "execute demo"), "DoSNMP")
			d.SubmitTask(task)
		}
		fmt.Println("tasks are submit")
		wg.Done()
	}()

	wg.Wait()
	fmt.Println("all tasks are finished")

}
