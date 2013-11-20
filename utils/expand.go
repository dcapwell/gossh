package utils

import (
  "strings"
  "strconv"
  "fmt"
)

const (
	LEFT  = "["
	RIGHT = "]"
	ELIPS = ".."
  COMM  = ","
)

/*
Expands the given string into 1 or more strings.  Expand has support for range, and list based expansions.  Range based expansions use the notation [start..end], where start is less than or equal to end.  List based expansions use a comma ',' to denote a new element.  Any number of range and list expansions are supported.
*/
func Expand(input string) ([]string, error) {
  /*
  Algo:
    split off COMM to get all list elements
    for each element, recursively expand all ranges
  */
  data := make([]string, 0)
  splits := expandSplit(input)
  for _, split := range splits {
    rang, err := expandRange(split)
    if err != nil {
      return nil, err
    }
    for _, r := range rang {
      data = append(data, r)
    }
  }
  return data, nil
}

// splits the input into multiple elements based off Expand's list rules
func expandSplit(input string) []string {
  return strings.Split(input, COMM)
}

// expand a given element based off Expand's range rules
func expandRange(input string) ([]string, error) {
  if leftBracket := strings.Index(input, LEFT); leftBracket >= 0 {
    if rightBracket := strings.Index(input, RIGHT); rightBracket > leftBracket {
      // get range param
      r := input[leftBracket + 1 : rightBracket]
      rang, err := explodeRangeParam(r)
      if err != nil {
        return nil, err
      }
      tail, err := expandRange(input[rightBracket + 1:])
      if err != nil {
        return nil, err
      }
      // merge ranges together
      size := len(rang) * len(tail)
      data := make([]string, size)
      count := 0
      for i := 0; i < len(rang); i++ {
        left := input[0:leftBracket] + strconv.Itoa(rang[i])
        for j := 0; j < len(tail); j++ {
          data[count] = left + tail[j]
          count++
        }
      }
      return data, nil
    }
  }
  return []string{input}, nil
}

func explodeRangeParam(input string) ([]int, error) {
  sides := strings.Split(input, ELIPS)
  if len(sides) != 2 {
    return nil, fmt.Errorf("Unable to parse range param %s.\n", input)
  }
  start, err := strconv.Atoi(sides[0])
  if err != nil {
    return nil, err
  }
  end, err := strconv.Atoi(sides[1])
  if err != nil {
    return nil, err
  }
  //TODO support this.  Its valid to want range to go in reverse
  if start > end {
    return nil, fmt.Errorf("Unable to parse %s; start value in range is larger than end.\n", input)
  }
  size := end - start + 1
  nums := make([]int, size)
  for i := 0; i < size; i++ {
    nums[i] = start + i
  }
  return nums, nil
}
