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
