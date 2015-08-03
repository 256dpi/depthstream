package main
import (
  "github.com/velovix/go-freenect"
  "fmt"
  "time"
)

func CountDevices() int {
  ctx, e1 := freenect.NewContext()
  defer ctx.Destroy()

  if e1 != nil {
    panic(e1)
  }

  cnt, e2 := ctx.DeviceCount()

  if e2 != nil {
    panic(e2)
  }

  return cnt;
}

func StreamDepthImage(port int, device int) {
  ctx, e1 := freenect.NewContext()
  defer ctx.Destroy()

  if e1 != nil {
    panic(e1)
  }

  ctx.SetLogLevel(freenect.LogDebug)

  kinect, e2 := ctx.OpenDevice(device)
  defer kinect.Destroy()

  if e2 != nil {
    panic(e2)
  }

  kinect.SetDepthCallback(func(device *freenect.Device, depth []uint16, timestamp uint32){
    fmt.Println(len(depth))
  })

  e3 := kinect.StartDepthStream(freenect.ResolutionMedium, freenect.DepthFormatMM)
  defer kinect.StopDepthStream()

  if e3 != nil {
    panic(e3)
  }

  start := time.Now()

  for time.Since(start).Seconds() < 10.0 {
    // Process freenect events
    err := ctx.ProcessEvents(0)
    if err != nil {
      fmt.Println(err)
      break
    }
  }
}
