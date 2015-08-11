package main

import (
  "encoding/binary"
  "github.com/cznic/sortutil"
)

func isPowerOfTwo(x int) bool {
  return (x & (x - 1)) == 0;
}

func reduce(data []uint16, n int) []uint16 {
  set := make([]uint16, len(data) / n / 2)

  for y := 0; y < 480 / n; y++ {
    for x := 0; x < 640 / n; x++ {
      set[y * 640 / n + x] = data[y * n * 640 + x * n];
    }
  }

  return set
}

func averagePixelBlock(data []uint16, x int, y int, blockSize int) uint16 {
  block := make(sortutil.Uint16Slice, 0, (blockSize + 1) * blockSize)

  for yy := y - blockSize; yy <= y + blockSize; yy++ {
    for xx := x - blockSize; xx <= x + blockSize; xx++ {
      if yy >= 0 && yy < 480 && xx >= 0 && xx < 640 {
        p := data[yy * 640 + xx]

        if p > 0 {
          block = append(block, p)
        }
      }
    }
  }

  if len(block) > 0 {
    block.Sort()
    return block[len(block) / 2]
  } else {
    return 0
  }
}

func interpolate(data []uint16, blockSize int) []uint16 {
  for y := 0; y < 480; y++ {
    for x := 0; x < 640; x++ {
      i := y * 640 + x

      if data[i] == 0 {
        data[i] = averagePixelBlock(data, x, y, blockSize)
      }
    }
  }

  return data
}

func Convert(c *Config, data []uint16) []byte {
  if c.interpolate > 0 {
    data = interpolate(data, c.interpolate)
  }

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
