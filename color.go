package main

func reduceColor(data []byte, n int) []byte {
  set := make([]byte, len(data) / n / n)

  for y := 0; y < 480 / n; y++ {
    for x := 0; x < 640 / n; x++ {
      p1 := (y * 640 / n + x) * 3
      p2 := (y * n * 640 + x * n) * 3

      set[p1] = data[p2]
      set[p1 + 1] = data[p2 + 1]
      set[p1 + 2] = data[p2 + 2]
    }
  }

  return set
}

func ConvertColor(c *Config, data []byte) []byte {
  if c.reduce > 0 && isPowerOfTwo(c.reduce) {
    return reduceColor(data, c.reduce)
  }

  return data
}
