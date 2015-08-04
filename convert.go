package main

import (
  "encoding/binary"
)

func isPowerOfTwo(x int) bool {
  return (x & (x - 1)) == 0;
}

func reduce(data []uint16, n int) []uint16 {
  set := make([]uint16, len(data) / n / 2)

  for x := 0; x < 640 / n; x++ {
    for y := 0; y < 480 / n; y++ {
      set[y * 640 / n + x] = data[y * n * 640 + x * n];
    }
  }

  return set
}

func Convert(c *Config, data []uint16) []byte {
  if c.reduce > 0 && isPowerOfTwo(c.reduce) {
    data = reduce(data, c.reduce)
  }

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
