package prototype

import (
  "testing"
  "time"
  "fmt"
  "log"
)

const (
  SUCCESS = iota
  FAILURE
)

type Task interface {
  Run() (interface{}, error)
}

type NoOP struct {}

func (n NoOP) Run() (interface{}, error) {
  return "no-op", nil
}

type PanicAtTheDisco struct {}

func (n PanicAtTheDisco) Run() (interface{}, error) {
  panic("at the disco")
}

type TaskResult struct {
  Duration  time.Duration
  Status    int
  Result    interface{}
}

type FailedTaskResult struct {
  Result    interface{}
  Error     error
}

func safeTaskRunner(task Task) (output interface{}, err error) {
  defer func() {
    if e := recover(); e != nil {
      log.Printf("task (%v) failed: %v", task, e)
      _, ok := e.(error)
      if !ok {
        err = fmt.Errorf("task: %v", e)
      }
    }
  }()
  return task.Run()
}

// executes Task.Run concurrently in size number of workers
// expects that the ch chan will be closed after data is finished
// if data will never finish and tasks will constantly be added, then its fine to leave
// ch open.  Pool will be blocked waiting for tasks to run
func Pool(ch chan Task, size int) chan TaskResult {
  results := make(chan TaskResult)
  // used like a countdown latch.  After size events, its safe to close results
  done := make(chan int)
  // create size number of workers to run tasks
  for i := 0; i < size; i++ {
    go func() {
      var start time.Time
      for task := range ch {
        // task could panic, need to handle this
        start = time.Now()
        // result, _ := task.Run()
        result, err := safeTaskRunner(task)
        if err == nil {
          results <- TaskResult{time.Since(start), SUCCESS, result}
        } else {
          results <- TaskResult{time.Since(start), FAILURE, FailedTaskResult{result, err}}
        }
      }
      // input is done, so close the results
      // log.Print("Closing results")
      // close(results)
      // close isn't safe since multiple functions will try to call it.
      // treat done as a countdown latch
      done <- 1
    }()
  }
  //TODO is there a cleaner way to close result?
  go func() {
    for i := 0; i < size; i++ {
      <-done
    }
    // if I don't close, then I can't use range on the channel
    close(results)
  }()
  return results
}

func TestPool(t *testing.T) {
  ch := make(chan Task, 10)
  results := Pool(ch, 5)
  go func() {
    for i := 0; i < 40; i++ {
      //log.Print("Adding NoOP task")
      ch <- NoOP{}
    }
    close(ch)
  }()

  for result := range results {
    log.Printf("Got back result %v\n", result)
  }
}

func TestPoolWithPanic(t *testing.T) {
  ch := make(chan Task, 10)
  results := Pool(ch, 5)
  go func() {
    for i := 0; i < 40; i++ {
      //log.Print("Adding PanicAtTheDisco task")
      ch <- PanicAtTheDisco{}
    }
    close(ch)
  }()

  for result := range results {
    log.Printf("Got back result %v\n", result)
  }
}
