package test

import (
	"runtime"
	"time"
)

func doTask() {
	//耗时炒作(模拟)
	time.Sleep(1 * time.Second)
}

//这里模拟的http接口,每次请求抽象为一个job
func handle() {
	for i := 0; i < 10; i++ {
		job := Job{}
		JobQueue <- job
	}
}

var (
	MaxWorker = runtime.NumCPU() //CPU核心数是最大并行数,worker多余并加快不了任务的执行速度,反而会增加切换的负担
	MaxQueue  = 200000
)

type Worker struct {
	quit chan bool
}

func NewWorker() Worker {
	return Worker{
		quit: make(chan bool)}
}

// Start method starts the run loop for the worker, listening for a quit channel in
// case we need to stop it
func (w Worker) Start() {
	go func() {
		for {
			select {
			case <-JobQueue:
				// we have received a work request.
				doTask()
			case <-w.quit:
				// we have received a signal to stop
				return
			}
		}
	}()
}

// Stop signals the worker to stop listening for work requests.
func (w Worker) Stop() {
	go func() {
		w.quit <- true
	}()
}

type Job struct {
}

var JobQueue = make(chan Job, MaxQueue)

type Dispatcher struct {
}

func NewDispatcher() *Dispatcher {
	return &Dispatcher{}
}

func (d *Dispatcher) Run() {
	// starting n number of workers
	for i := 0; i < MaxWorker; i++ {
		worker := NewWorker()
		worker.Start()
	}
}
