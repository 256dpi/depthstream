package main

import (
  "encoding/binary"
  "github.com/cznic/sortutil"
)

func reduceDepth(data []uint16, n int) []uint16 {
  set := make([]uint16, len(data) / n / n)

  for y := 0; y < 480 / n; y++ {
    for x := 0; x < 640 / n; x++ {
      set[y * 640 / n + x] = data[y * n * 640 + x * n];
    }
  }

  return set
}

func averagePixelBlock(data []uint16, x int, y int, w int, h int, blockSize int) uint16 {
  block := make(sortutil.Uint16Slice, 0, (blockSize + 1) * blockSize)

  for yy := y - blockSize; yy <= y + blockSize; yy++ {
    for xx := x - blockSize; xx <= x + blockSize; xx++ {
      if yy >= 0 && yy < h && xx >= 0 && xx < w {
        p := data[yy * w + xx]

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

func interpolateDepth(data []uint16, w int, h int, blockSize int) []uint16 {
  for y := 0; y < h; y++ {
    for x := 0; x < w; x++ {
      i := y * w + x

      if data[i] == 0 {
        data[i] = averagePixelBlock(data, x, y, w, h, blockSize)
      }
    }
  }

  return data
}

func ConvertDepth(c *Config, data []uint16) []byte {
  var w, h int

  if c.reduce > 0 && isPowerOfTwo(c.reduce) {
    data = reduceDepth(data, c.reduce)
    w = 640 / c.reduce
    h = 480 / c.reduce
  } else {
    w = 640
    h = 480
  }

  if c.interpolate > 0 {
    data = interpolateDepth(data, w, h, c.interpolate)
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
