package prototype

import (
  "testing"
  "strings"
  "log"
  "fmt"
  "strconv"
)

const (
  LEFT = "["
  RIGHT = "]"
  ELIPS = ".."
)

func elips(elipsStr string) ([]int, error) {
  sides := strings.Split(elipsStr, ELIPS)
  if len(sides) != 2 {
    return nil, fmt.Errorf("Unable to parse %s\n", elipsStr)
  }
  start, err := strconv.Atoi(sides[0])
  if err != nil {
    return nil, err
  }
  end, err := strconv.Atoi(sides[1])
  if err != nil {
    return nil, err
  }
  if start > end {
    return nil, fmt.Errorf("Unable to parse %s; end range larger than start\n", elipsStr)
  }
  size := end - start + 1
  nums := make([]int, size)
  for i := 0; i < size; i++ {
    nums[i] = start + i
  }
  return nums, nil
}

func expand(str string) ([]string, error) {
  leftBracket := strings.Index(str, LEFT)
  if leftBracket >= 0 {
    rightBracket := strings.Index(str, RIGHT)
    if rightBracket > leftBracket {
      r := str[leftBracket+1:rightBracket]
      rang, err := elips(r)
      if err != nil {
        return nil, err
      }
      tail, err := expand(str[rightBracket + 1:])
      if err != nil {
        return nil, err
      }
      size := len(rang) * len(tail)
      data := make([]string, size)
      count := 0
      for i := 0; i < len(rang); i++ {
        left := str[0:leftBracket] + strconv.Itoa(rang[i])
        for j := 0; j < len(tail); j++ {
          data[count] = left + tail[j]
          count++
        }
      }
      return data, nil
    }
  }
  return []string{str}, nil
}

func TestTwoHostExpantion(t *testing.T) {
  hostStr := "example[1..2]"
  hosts, err := expand(hostStr)
  if err != nil {
    t.Errorf("Error while parsing host: %v\n", err)
    t.FailNow()
  }

  host := hosts[0]
  log.Printf("Host %s\n", host)
  if host != "example1" {
    t.Errorf("Expected example1, but found %s\n", host)
    t.FailNow()
  }

  host = hosts[1]
  log.Printf("Host %s\n", host)
  if host != "example2" {
    t.Errorf("Expected example2, but found %s\n", host)
    t.FailNow()
  }
}

func TestTwoHostExpantionInMiddle(t *testing.T) {
  hostStr := "example[1..2].ic"
  hosts, err := expand(hostStr)
  if err != nil {
    t.Errorf("Error while parsing host: %v\n", err)
    t.FailNow()
  }

  host := hosts[0]
  log.Printf("Host %s\n", host)
  if host != "example1.ic" {
    t.Errorf("Expected example1.ic, but found %s\n", host)
    t.FailNow()
  }

  host = hosts[1]
  log.Printf("Host %s\n", host)
  if host != "example2.ic" {
    t.Errorf("Expected example2.ic, but found %s\n", host)
    t.FailNow()
  }
}

func TestMultiExpantions(t *testing.T) {
  hostStr := "example[1..2].ic[4..5].com"
  hosts, err := expand(hostStr)
  if err != nil {
    t.Errorf("Error while parsing host: %v\n", err)
    t.FailNow()
  }

  if len(hosts) != 4 {
    t.Errorf("Number of expantions not matched: %s", hosts)
  }

  log.Printf("Hosts %s\n", hosts)

  host := hosts[0]
  log.Printf("Host %s\n", host)
  if host != "example1.ic4.com" {
    t.Errorf("Expected example1.ic4.com, but found %s\n", host)
    t.FailNow()
  }

  host = hosts[1]
  log.Printf("Host %s\n", host)
  if host != "example1.ic5.com" {
    t.Errorf("Expected example1.ic5.com, but found %s\n", host)
    t.FailNow()
  }

  host = hosts[2]
  log.Printf("Host %s\n", host)
  if host != "example2.ic4.com" {
    t.Errorf("Expected example2.ic4.com, but found %s\n", host)
    t.FailNow()
  }

  host = hosts[3]
  log.Printf("Host %s\n", host)
  if host != "example2.ic5.com" {
    t.Errorf("Expected example2.ic5.com, but found %s\n", host)
    t.FailNow()
  }
}
