package main

import (
  "encoding/binary"
)

func Convert(c *Config, data []uint16) []byte {
  buf := make([]byte, len(data) * 2)

  for i, p := range data {
    if c.bigendian {
      binary.BigEndian.PutUint16(buf[i*2:], p)
    } else {
      binary.LittleEndian.PutUint16(buf[i*2:], p)
    }
  }

  return buf
}
