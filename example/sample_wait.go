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

var wg sync.WaitGroup

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	d := npd.NewDispatcherWithWait(30, 30, &wg)

	d.Run()
	defer d.Stop()

	wg.Add(1)

	go func() {
		for i := 0; i < 30; i++ {
			task := npd.MakeTask(npd.TASK_NORMAL, NewMyTask(i, "execute demo"), "DoSNMP")
			d.SubmitTask(task)
		}
		fmt.Println("Tasks sent.")
		wg.Done()
	}()

	wg.Wait()
	fmt.Println("All tasks are done")

	select {}

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
