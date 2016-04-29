package npd

import (
	"reflect"
	"sync"
)

//执行器结构
type Executor struct {
	TaskPool chan chan Task //任务池
	TaskChan chan Task      //任务通道
	wg       *sync.WaitGroup
	quit     chan bool
	wait     bool
	idle     bool //是否空闲
}

//构建执行器
func NewExecutor(taskPool chan chan Task) Executor {
	return Executor{
		TaskPool: taskPool,
		TaskChan: make(chan Task),
		wait:     false,
		quit:     make(chan bool),
		idle:     false,
	}
}

//构建带等待的执行器
func NewExecutorWithWait(taskPool chan chan Task, wg *sync.WaitGroup) Executor {
	t := NewExecutor(taskPool)
	t.wg = wg
	t.wait = true
	return t
}

//开启执行模式
func (e *Executor) Start() {
	go func() {
		for {
			e.TaskPool <- e.TaskChan
			select {
			case task := <-e.TaskChan:
				e.idle = false
				if task.Type == TASK_NORMAL {
					reflect.ValueOf(*task.TargetObj).MethodByName(task.TargetFunc).Call([]reflect.Value{})
				}
				e.idle = true
				if e.wait {
					e.wg.Done()
				}
			case <-e.quit:
				println("executor quit")
				return
			}
		}
	}()
}

//停止执行模式(不是立刻停止,而是发送停止信号呼叫当前任务停止)
func (e *Executor) Stop() bool {
	if e.idle {
		e.quit <- true
		return true
	}
	return false
}
