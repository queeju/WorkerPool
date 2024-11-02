package main

import (
	"fmt"
	"log"

	"github.com/queeju/WorkerPool/workerpool"
)

func main() {
	var numOfWorkers uint = 3

	wp := workerpool.NewWorkerPool(numOfWorkers)

	var jCount int

	addJobs(10, &jCount, wp)

	err := wp.AddWorker()
	if err != nil {
		log.Printf("Error adding worker: %s", err)
		return
	}

	addJobs(10, &jCount, wp)

	for i := 0; i < 3; i++ {
		if err := wp.RemoveWorker(); err != nil {
			log.Printf("Error removing worker: %s", err)
		}
	}

	addJobs(10, &jCount, wp)

	if err := wp.Stop(); err != nil {
		log.Printf("Error stopping worker pool: %s", err)
	}
}

func addJobs(num int, jCount *int, wp *workerpool.WorkerPool) {
	for i := 0; i < num; i++ {
		if err := wp.AddJob(fmt.Sprintf("Job #%d", *jCount)); err != nil {
			log.Printf("Error adding job: %s", err)
			return
		}
		*jCount++
	}
}
