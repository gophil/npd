package main

import (
	"flag"
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
	npd.GetLogger().Infoln(m.number, "==>", m.message)
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

	//创建分发器
	d := npd.NewDispatcherWithWait(num, num, &wg)

	d.Run()
	defer d.Stop()

	wg.Add(1)

	go func() {
		for i := 0; i < 600; i++ {
			task := npd.CreateTask(NewMyTask(i, "execute demo"), "DoSNMP")
			d.SubmitTask(task)
		}
		npd.GetLogger().Infoln("tasks are submit")
		wg.Done()
	}()

	wg.Wait()
	npd.GetLogger().Infoln("all tasks are finished")

}
