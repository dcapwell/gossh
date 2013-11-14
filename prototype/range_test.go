package prototype

import (
  "testing"
  "time"
  "log"
)

type Range struct {
  Start int32
  End   int32
}

func (r *Range) Iter() <-chan int32 {
  // channel only allows one, so ch <- i will block till something pops, so only generate data when needed
  ch := make(chan int32, 1)
  go func() {
    for i := r.Start; i <= r.End; i++ {
      ch <- i
    }
    // close channel so range will finish
    close(ch)
  }()
  return ch
}

func TestRange(t *testing.T) {
  r := &Range{Start: 1, End: 32}

  for i := range r.Iter() {
    log.Printf("Number: %d\n", i)
  }
}

func TestLazyRange(t *testing.T) {
  r := &Range{Start: 1, End: 10}

  for i := range r.Iter() {
    log.Printf("Number: %d\n", i)
    time.Sleep(1 * time.Second)
  }
}
