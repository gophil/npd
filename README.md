## npd

一个简单的自定义任务分发模块


`安装第三方依赖`

*   > go get github.com/cihub/seelog



###如何使用 ?

`Example: `

1.执行一个简单的任务

```go
    
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

    func (m *MyTask) DoSNMP() {
        fmt.Println(m.message, " -> ", m.number)
    }

    func main() {
        runtime.GOMAXPROCS(runtime.NumCPU())

        d := npd.NewDispatcher(30, 30)

        d.Run()
        defer d.Stop()

        for i := 0; i < 30; i++ {
            task := npd.MakeTask(npd.TASK_NORMAL, NewMyTask(i, "execute demo"), "DoSNMP")
            d.SubmitTask(task)
        }

        select {}

    }

  


```


2.执行一个带等待的任务

```go
    
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
            fmt.Println("tasks are submit ")
            wg.Done()
        }()

        wg.Wait()
        fmt.Println("all tasks are finished")

    }


```

3.执行一个支持消息发送的任务
```go
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
	Number  int
	Message string
}

func NewMyTask(number int, message string) *MyTask {
	return &MyTask{
		Number:  number,
		Message: message,
	}
}

func (m *MyTask) DoSNMP() {
	time.Sleep(500 * time.Millisecond)
	fmt.Println(m.Message, " -> ", m.Number)
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

	//设置消息发送函数
	d.SetMF(func(task npd.Task) {
		no := (*task.TargetObj).(*MyTask)
		println("处理数据上报:", no.Number)
	})

	d.Run()
	defer d.Stop()

	wg.Add(1)
	mpwg.Add(1)

	go func() {
		for i := 0; i < 530; i++ {
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

```


####自定义任务结构与方法说明

```go
//任务对象需要struct类型, 不支持interface类型
type Custom struct {
    
}
```


####结构体接收方法为无入参的方法

eg:

```go

func (c *Custom) foo() {
    
    ...
}

```

####方法只支持无返回值和错误类型返回值的形式

eg:

```go

func (c *Custom) foo2() {
    
    ...
}

func (c *Custom) foo2() error {
    
    ...
}

```
