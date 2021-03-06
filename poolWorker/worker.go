package poolWorker

import (
	log "github.com/sirupsen/logrus"
	"runtime"
	"sync"
)

type Work interface {
	Do()
}

type worker struct {
	id        uint
	done      *sync.WaitGroup
	readyPool chan chan Work //get work from the boss
	work      chan Work
	quit      chan bool
}

func NewWorker(id uint, readyPool chan chan Work, done *sync.WaitGroup) *worker {
	return &worker{
		id:        id,
		done:      done,
		readyPool: readyPool,
		work:      make(chan Work),
		quit:      make(chan bool),
	}
}

func (w *worker) Process(work Work) {
	//DoOnCreate the work
	defer func() { //Capture any panic or error
		if r := recover(); r != nil {
			const size = 64 << 10
			buf := make([]byte, size)
			buf = buf[:runtime.Stack(buf, false)]
			log.Error("panic running process: %v\n%s\n", r, buf)
		}
	}()
	work.Do()
}

func (w *worker) Start() {
	go func() {
		w.done.Add(1) // wait for me
		for {
			w.readyPool <- w.work //hey i am ready to work on new job
			select {
			case work := <-w.work: // hey i am waiting for new job
				w.Process(work) // ok i am on it
			case <-w.quit:
				w.done.Done() // ok i am here i finished my all jobs
				return
			}
		}
	}()
}

func (w *worker) Stop() {
	//tell worker to stop after current process
	w.quit <- true
}
