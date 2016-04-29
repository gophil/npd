package npd

import (
	"sync"
	"time"
)

//任务分发器
type Dispatcher struct {
	maxExecutors    int
	queueBufferSize int
	wg              *sync.WaitGroup
	wait            bool
	taskQueue       chan Task
	taskPool        chan chan Task
	executors       []Executor
	quit            chan bool
	limiter         <-chan time.Time
}

//创建分发器
func NewDispatcher(maxExecutors, queueBufferSize int) *Dispatcher {
	dispatcher := &Dispatcher{
		maxExecutors:    maxExecutors,
		queueBufferSize: queueBufferSize,
		taskPool:        make(chan chan Task, maxExecutors),
		wait:            false,
		quit:            make(chan bool),
		executors:       make([]Executor, maxExecutors),
	}

	if queueBufferSize != 0 {
		dispatcher.taskQueue = make(chan Task, maxExecutors)
	} else {
		dispatcher.taskQueue = make(chan Task)
	}
	return dispatcher
}

//创建带等待的分发器
func NewDispatcherWithWait(maxExecutors, queueBufferSize int, wg *sync.WaitGroup) *Dispatcher {
	dispatcher := NewDispatcher(maxExecutors, queueBufferSize)
	dispatcher.wg = wg
	dispatcher.wait = true
	return dispatcher
}

func (dispatcher *Dispatcher) SubmitTask(task Task) {
	if dispatcher.wait {
		dispatcher.wg.Add(1)
	}
	dispatcher.taskQueue <- task
}

func (dispatcher *Dispatcher) Run() {

	for i := 0; i < dispatcher.maxExecutors; i++ {
		if dispatcher.wait {
			dispatcher.executors[i] = NewExecutorWithWait(dispatcher.taskPool, dispatcher.wg)
		} else {
			dispatcher.executors[i] = NewExecutor(dispatcher.taskPool)
		}
		//开启执行
		dispatcher.executors[i].Start()
	}

	//开启调度
	go dispatcher.dispatch()
}

func (dispatcher *Dispatcher) RunWithLimiter(limiterGap time.Duration) {
	dispatcher.limiter = time.Tick(limiterGap)
	dispatcher.Run()
}

//分发处理
func (dispatcher *Dispatcher) dispatch() {
	defer dispatcher.shutdown()
	for {
		select {
		case task := <-dispatcher.taskQueue:

			executorTaskChan := <-dispatcher.taskPool

			if dispatcher.limiter != nil {
				<-dispatcher.limiter
			}

			executorTaskChan <- task

		case <-dispatcher.quit:
			return
		}
	}
}

func (dispatcher *Dispatcher) shutdown() {
	for _, e := range dispatcher.executors {
		for !e.Stop() { //一直处理停止
		}
	}
	close(dispatcher.taskPool)
	close(dispatcher.taskQueue)
}

//停止开关
func (dispatcher *Dispatcher) Stop() {
	dispatcher.quit <- true
}
