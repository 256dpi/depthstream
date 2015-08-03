package main

import (
  "github.com/docopt/docopt-go"
  "strconv"
)

type Config struct {
  info bool
  start bool
  device int
  port int
}

func ParseConfig() *Config {
  usage := `Depthstream.

Usage:
    depthstream info
    depthstream start <device> <port>`

  a, _ := docopt.Parse(usage, nil, true, "", false)

  return &Config{
    getBool(a["info"]),
    getBool(a["start"]),
    getInt(a["<device>"]),
    getInt(a["<port>"]),
  }
}

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
