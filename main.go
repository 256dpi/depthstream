package main

import (
  "fmt"
  "syscall"
  "os/signal"
  "os"
)

func main() {
  c := ParseConfig()

  if c.info {
    count := CountDevices()

    if(count == 1) {
      fmt.Printf("There is one Kinect connected.\n")
    } else {
      fmt.Printf("There are %d Kinects connected.\n", count)
    }
  } else if c.start {
    if c.device >= 0 && c.port > 100 {
      fmt.Printf("Start stream server on port %d with data from kinect %d...\n", c.port, c.device)

      data := make(chan []uint16)
      queue := make(chan *connection)
      stream := make(chan *connection)

      relay := NewRelay(queue, stream)
      relay.Start(c.port)

      depthStream := NewDepthStream(data)
      depthStream.Open(0)

      go func(){
        var cache []uint16
        var list []*connection
        for {
          select {
          case cache = <-data:
            for _, conn := range list {
              conn.send <- Convert(cache)
            }
          case conn := <-queue:
            conn.send <- Convert(cache)
          case conn := <-stream:
            list = append(list, conn)
          }
        }
      }()

      finish := make(chan os.Signal, 1)
      signal.Notify(finish, syscall.SIGINT, syscall.SIGTERM)

      <-finish

      depthStream.Close()
      relay.Stop()
    } else {
      fmt.Printf("Specify a device id >= 0 and port >= 100!\n")
    }
  }
}
