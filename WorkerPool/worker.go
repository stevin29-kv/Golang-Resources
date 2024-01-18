package main

import (
	"context"
	"log"
)

const (
	defaultWorkerSize       = 1
	SQSWorkerSize           = 3
	PaymentCancelWorkerSize = 3
)

type WorkerPool interface {
	Run(ctx context.Context)
	AddTask(task func())
}

type workerPool struct {
	maxWorker    int
	queuedTaskCh chan func()
}

// NewWorkerPool will create an instance of WorkerPool.
// How to use:
//
//	workerSize := 10
//	wp := workerpool.NewWorkerPool(workerSize)
//	wp.Run(ctx)
//
//	 for i := 0; i < totalTask; i++ {
//	  id := i + 1
//	  wp.AddTask(func() {
//	  logger.WithContext(ctx).Infof("processing task %d", i)
//	 })
//	}

func NewWorkerPool(maxWorker int) WorkerPool {
	if maxWorker == 0 {
		maxWorker = defaultWorkerSize
	}
	wp := &workerPool{
		maxWorker:    maxWorker,
		queuedTaskCh: make(chan func(), maxWorker),
	}
	return wp
}

func (wp *workerPool) Run(ctx context.Context) {
	wp.run(ctx)
}

func (wp *workerPool) AddTask(task func()) {
	wp.queuedTaskCh <- task
}

func (wp *workerPool) GetTotalQueuedTask() int {
	return len(wp.queuedTaskCh)
}

func (wp *workerPool) run(ctx context.Context) {
	for i := 0; i < wp.maxWorker; i++ {
		wID := i + 1
		log.Printf("[WorkerPool] Worker %d has been spawned", wID)

		go func(workerID int) {
			for task := range wp.queuedTaskCh {
				// log.Printf("[WorkerPool] Worker %d started processing task", wID)
				task()
				// log.Printf("[WorkerPool] Worker %d finished processing task", wID)
			}
		}(wID)
	}
}
