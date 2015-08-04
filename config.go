package main

import (
  "github.com/docopt/docopt-go"
  "strconv"
)

type Config struct {
  info bool
  device int
  port int
}

func ParseConfig() *Config {
  usage := `Depthstream.

Usage:
  depthstream [options]

Options:
  -h --help         Show this screen.
  -i --info         Show connected Kinects.
  -p --port=<num>   Port for server. [default: 9090].
  -d --dev=<id>     Device to open. [default: 0].
`

  a, _ := docopt.Parse(usage, nil, true, "", false)

  return &Config{
    info: getBool(a["--info"]),
    device: getInt(a["--dev"]),
    port: getInt(a["--port"]),
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
