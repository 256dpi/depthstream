package main

import (
  "github.com/docopt/docopt-go"
  "strconv"
)

type Config struct {
  info bool
  device int
  port int
  bigendian bool
}

func ParseConfig() *Config {
  usage := `Depthstream.

Usage:
  depthstream [options]

Options:
  -h --help             Show this screen.
  -i --info             Show connected Kinects.
  -p --port=<n>         Port for server. [default: 9090].
  -d --device=<n>       Device to open. [default: 0].
  -b --bigendian        Use big endian encoding.
`

  a, _ := docopt.Parse(usage, nil, true, "", false)

  return &Config{
    info: getBool(a["--info"]),
    device: getInt(a["--dev"]),
    port: getInt(a["--port"]),
    bigendian:getBool(a["--bigendian"]),
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
