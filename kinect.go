package main

import (
  "github.com/velovix/go-freenect"
  "fmt"
  "sync"
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

type DepthStream struct {
  data chan []uint16
  stop chan int
  wg *sync.WaitGroup
  ctx *freenect.Context
  dev *freenect.Device
}

func NewDepthStream(ch chan []uint16) *DepthStream {
  return &DepthStream{
    ch,
    make(chan int),
    &sync.WaitGroup{},
    nil,
    nil,
  }
}

func process(ds *DepthStream){
  defer ds.wg.Done()

  for {
    select {
    case _ = <-ds.stop:
      return
    default:
      ds.ctx.ProcessEvents(0);
    }
  }
}

func (ds *DepthStream) Open(device int) {
  var (
    err error
    ctx freenect.Context
    dev freenect.Device
  )

  ctx, err = freenect.NewContext()
  ds.ctx = &ctx

  if err != nil {
    panic(err)
  }

  dev, err = ds.ctx.OpenDevice(device)
  ds.dev = &dev

  if err != nil {
    panic(err)
  }

  ds.dev.SetDepthCallback(func(device *freenect.Device, depth []uint16, timestamp uint32){
    ds.data <- depth
  })

  err = ds.dev.StartDepthStream(freenect.ResolutionMedium, freenect.DepthFormatMM)

  if err != nil {
    panic(err)
  }

  ds.wg.Add(1)
  go process(ds)
}

func (ds *DepthStream) Close() {
  close(ds.stop)
  ds.wg.Wait()

  fmt.Println("Clean up...")
  ds.dev.StopDepthStream()
  ds.dev.Destroy()
  ds.ctx.Destroy()
}
