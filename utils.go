package main

import "strconv"

func getBool(field interface{}) bool {
  if bol, ok := field.(bool); ok {
    return bol
  } else {
    return false
  }
}

func getString(field interface{}) string {
  if str, ok := field.(string); ok {
    return str;
  } else {
    return ""
  }
}

func getInt(field interface{}) int {
  if num, ok := strconv.Atoi(getString(field)); ok == nil {
    return num
  } else {
    return 0
  }
}

func isPowerOfTwo(x int) bool {
  return (x & (x - 1)) == 0;
}
