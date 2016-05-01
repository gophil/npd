package main

import (
	"flag"
	"fmt"
	"github.com/gophil/npd"
	_ "net/http/pprof"
	"runtime"
	"strconv"
	"sync"
	"time"
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
	time.Sleep(500 * time.Millisecond)
	fmt.Println(m.message, " -> ", m.number)
}

var (
	work_num = flag.String("w", "100", "num of worker num") //执行的协程数量
)

func main() {

	runtime.GOMAXPROCS(runtime.NumCPU())

	flag.Parse()
	num, err := strconv.Atoi(*work_num)
	if err != nil {
		num = 100 //默认数
	}
	npd.GetLogger().Infof("executors size : %d\n", num)

	var wg sync.WaitGroup
	var mpwg sync.WaitGroup

	d := npd.NewDispatcherWithMQ(num, num, &wg, &mpwg)

	d.Run()
	defer d.Stop()

	wg.Add(1)
	mpwg.Add(1)

	go func() {
		for i := 0; i < 30; i++ {
			task := npd.CreateTask(NewMyTask(i, "execute demo"), "DoSNMP")
			d.SubmitTask(task)
		}
		fmt.Println("tasks are submit")
		wg.Done()
		mpwg.Done()
	}()

	wg.Wait()
	mpwg.Wait()
	fmt.Println("all tasks are finished")

}
