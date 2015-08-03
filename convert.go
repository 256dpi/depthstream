package main

func Convert(data []uint16) []byte {
  minData := make([]byte, len(data))

  for i, p := range data {
    minData[i] = byte(p / 40)
  }

  return minData
}
