package main

import (
  "encoding/binary"
)

func Convert(data []uint16) []byte {
  buf := make([]byte, len(data) * 2)

  for i, p := range data {
    // TODO: select endianness using a CLI option
    binary.LittleEndian.PutUint16(buf[i*2:], p)
  }

  return buf
}
