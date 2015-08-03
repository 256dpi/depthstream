package main

import (
  "fmt"
  "github.com/velovix/go-freenect"
  "os"
  "time"
)

var (
  context freenect.Context
  kinect freenect.Device
  frameCnt int
)

func onDepthFrame(device *freenect.Device, depth []uint16, timestamp uint32) {
  frameCnt++
}

func init() {
  context, e1 := freenect.NewContext()

  if e1 != nil {
    panic(e1)
  }

  context.SetLogLevel(freenect.LogDebug)

  count, e2 := context.DeviceCount()

  if e2 != nil {
    panic(e2)
  }

  if count == 0 {
    fmt.Println("could not find any devices")
    os.Exit(1)
  }

  kinect, e3 := context.OpenDevice(0)

  if e3 != nil {
    panic(e3)
  }

  kinect.SetDepthCallback(onDepthFrame)

  e4 := kinect.StartDepthStream(freenect.ResolutionMedium, freenect.DepthFormatMM)

  if e4 != nil {
    panic(e4)
  }
}

func main() {
  initTime := time.Now()

  for time.Since(initTime).Seconds() < 10.0 {
    e1 := context.ProcessEvents(0)

    if e1 != nil {
      panic(e1)
    }
  }

  fmt.Println("Processed", frameCnt, "frames in 10 seconds.")

  kinect.StopDepthStream()
  kinect.Destroy()
  context.Destroy()
}
