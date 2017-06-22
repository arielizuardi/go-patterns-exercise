package dispatcher

import (
	"fmt"
	"log"
)

// Communicating via this Queue
var JobQueue chan Job

type Job struct {
	DeviceToken    string
	NotificationID string
}

func (j *Job) Do() error {
	fmt.Printf(`Done sending notif %v to %v`+"\n", j.NotificationID, j.DeviceToken)
	return nil
}

type Dispatcher struct {
	WorkerPool chan chan Job
	maxWorkers int
}

func (d *Dispatcher) Run() {
	fmt.Printf(`Starting worker %v`+"\n", d.maxWorkers)
	for i := 0; i < d.maxWorkers; i++ {
		fmt.Printf(`Spawn worker #%v`+"\n", i+1)
		worker := NewWorker(d.WorkerPool)
		worker.Start()
	}

	go d.dispatch()
}

func NewDispatcher(maxWorkers int) *Dispatcher {
	pool := make(chan chan Job, maxWorkers)
	return &Dispatcher{pool, maxWorkers}
}

func (d *Dispatcher) dispatch() {

	for {
		select {
		case job := <-JobQueue:

			go func(job Job) {
				//obtain the available worker from workerpool
				//will be waiting until get one
				jobChannel := <-d.WorkerPool
				jobChannel <- job
			}(job)
		}
	}

}

// Worker represents the worker that executes the job
type Worker struct {
	WorkerPool chan chan Job
	JobChannel chan Job
	quit       chan bool
}

// Start method starts the run loop for the worker, listening for a quit channel in
// case we need to stop it
func (w Worker) Start() {
	go func() {
		for {
			// register the current worker into the worker queue.
			w.WorkerPool <- w.JobChannel
			select {
			case job := <-w.JobChannel:
				fmt.Printf(`Receiving Job Push to %v ...`+"\n", job.DeviceToken)
				// we have received a work request.
				if err := job.Do(); err != nil {
					log.Println("Error sending")
				}

			case <-w.quit:
				// we have received a signal to stop
				return
			}
		}
	}()
}

func (w Worker) Stop() {
	go func() {
		w.quit <- true
	}()
}

func NewWorker(workerPool chan chan Job) Worker {
	return Worker{
		WorkerPool: workerPool,
		JobChannel: make(chan Job),
		quit:       make(chan bool)}
}
