package main

import (
  "github.com/velovix/go-freenect"
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
  depth chan []uint16
  color chan []byte
  stop chan int
  wg *sync.WaitGroup
  ctx *freenect.Context
  dev *freenect.Device
  withColor bool
}

func NewDepthStream() *DepthStream {
  return &DepthStream{
    depth: make(chan []uint16),
    color: make(chan []byte),
    stop: make(chan int),
    wg: &sync.WaitGroup{},
    withColor: false,
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

func (ds *DepthStream) Open(device int, withColor bool) {
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

  if withColor {
    ds.withColor = true

    ds.dev.SetVideoCallback(func(device *freenect.Device, video []byte, timestamp uint32){
      ds.color <- video;
    })

    err = ds.dev.StartVideoStream(freenect.ResolutionMedium, freenect.VideoFormatRGB)

    if err != nil {
      panic(err)
    }
  }

  ds.dev.SetDepthCallback(func(device *freenect.Device, depth []uint16, timestamp uint32){
    ds.depth <- depth
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

  if ds.withColor {
    ds.dev.StopVideoStream()
  }

  ds.dev.StopDepthStream()
  ds.dev.Destroy()
  ds.ctx.Destroy()
}
