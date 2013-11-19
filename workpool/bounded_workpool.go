package workpool

import (
	"fmt"
	"log"
)

const (
	// Task completed without any errors or panics
	SUCCESS = iota
	// Task completed with an error or panic
	FAILURE
)

// Function to run within the work pool
type Task func() (interface{}, error)

/*
type Task interface {
  Run() (interface{}, error)
}
*/

// Result of the task run
type TaskResult struct {
	// Indicates if the result will contain task's output or an ErrorResult containing task's error output
	Status int
	// Task or ErrorResult output
	Result interface{}
}

// Wrapper for error's returned from the task
type ErrorResult struct {
	// Task's output, may be nil
	Result interface{}
	// Task's error or a wrapped panic
	Error error
}

// Create a new WorkPool with 1,000 max resources
func NewWorkPool() WorkPool {
	return &workPoolImpl{maxWorkers: 1000, currentWorkers: 0}
}

// Create a new WorkPool with max as the number of resources.  If max is not a positive value, then an error is returned
func NewWorkPoolWithMax(max int) (WorkPool, error) {
	if max <= 0 {
		return nil, fmt.Errorf("Unable to make new work pool with max value %d; max must be a positive value.\n", max)
	}
	return &workPoolImpl{maxWorkers: max, currentWorkers: 0}, nil
}

// A bounded pool of workers to run given tasks.  Currently the bounded max size is
// optimistic.  This means that the max number of workers may exceed Max if multiple
// users use the same pool concurrently.  This is only for the short term.  The bounding
// will become strong later when I learn golang locking/CAS.
type WorkPool interface {
	// Run the tasks in a new pool of min(Remaining, maxPoolSize).  If there are not more than
	// minPoolSize resources remaining, then an error will be returned
	Run(tasks chan Task, minPoolSize, maxPoolSize int) (chan TaskResult, error)

	// How many workers are free to use
	Remaining() int

	// Max number of works allowed at any given point
	Max() int

	// How many workers are currently in use
	Size() int
}

type workPoolImpl struct {
	maxWorkers     int
	currentWorkers int
}

func (w *workPoolImpl) Remaining() int {
	return w.maxWorkers - w.currentWorkers
}

func (w *workPoolImpl) Max() int {
	return w.maxWorkers
}

func (w *workPoolImpl) Size() int {
	return w.currentWorkers
}

func (w *workPoolImpl) Run(tasks chan Task, minPoolSize, maxPoolSize int) (chan TaskResult, error) {
	remaining := w.Remaining()
	poolSize := min(maxPoolSize, remaining)
	if poolSize < minPoolSize {
		return nil, fmt.Errorf("Unable to create new pool; remaining (%d) resources is less than min (%d) requested\n", remaining, minPoolSize)
	}
	w.currentWorkers += poolSize

	return w.newPool(tasks, poolSize), nil
}

func (w *workPoolImpl) newPool(tasks chan Task, size int) chan TaskResult {
	// build pool
	results := make(chan TaskResult)
	// used like a countdown latch.  After size events, its safe to close results
	done := make(chan int)
	// create size number of workers to run tasks
	for i := 0; i < size; i++ {
		go func() {
			// need to be in defer in case a panic propagates up
			defer func() {
				// treat done as a countdown latch
				done <- 1
			}()
			for task := range tasks {
				// result, _ := task.Run()
				result, err := safeTaskRunner(task)
				if err == nil {
					results <- TaskResult{SUCCESS, result}
				} else {
					results <- TaskResult{FAILURE, ErrorResult{result, err}}
				}
			}
		}()
	}
	//TODO is there a cleaner way to close result?
	go func() {
		for i := 0; i < size; i++ {
			<-done
			// worker resource is done, so add back to pool
			w.currentWorkers -= 1
		}
		// if I don't close, then I can't use range on the channel
		close(results)
	}()
	return results
}

// math.Min takes float64, rather not convert back and forth for something so simple
func min(left, right int) int {
	if left <= right {
		return left
	}
	return right
}

// task runner that will recover from any panics, and return it as an error
func safeTaskRunner(task Task) (output interface{}, err error) {
	// recover from any panic and return panic as an error
	defer func() {
		if e := recover(); e != nil {
			log.Printf("task (%v) failed: %v", task, e)
			_, ok := e.(error)
			if !ok {
				err = fmt.Errorf("task: %v", e)
			}
		}
	}()
	return task()
}
