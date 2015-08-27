package main

import (
  "fmt"
  "syscall"
  "os/signal"
  "os"
  "time"
)

func start(c *Config) {
  relay := NewRelay()
  relay.Start(c.port)

  fmt.Printf("Server launched on port %d!\n", c.port)

  depthStream := NewDepthStream()

  go func(){
    var depthCache []byte
    var colorCache []byte
    var cache []byte

    i := 0
    list := make(map[*connection]bool)
    for {
      select {
      case depth := <-depthStream.depth:
        if(c.skip > 0) {
          i++;

          if(i > c.skip) {
            i = 0
          } else {
            break
          }
        }

        depthCache = Convert(c, depth)

        cache = make([]byte, len(depthCache) + len(colorCache))
        copy(cache[:len(depthCache)], depthCache)
        copy(cache[len(depthCache):], colorCache)

        for conn, _ := range list {
          conn.send <- cache
        }
      case color := <-depthStream.color:
        colorCache = color
      case conn := <-relay.queue:
        conn.send <- cache
      case conn := <-relay.stream:
        list[conn] = true
      case conn := <-relay.unstream:
        if _, ok := list[conn]; ok {
          delete(list, conn)
        }
      }
    }
  }()

  depthStream.Open(c.device, c.color)

  ticker := time.NewTicker(1 * time.Second)

  go func() {
    for {
      select {
      case <-ticker.C:
        fmt.Printf("\033[2K\033[1GClients: %d", len(relay.connections))
      }
    }
  }()

  finish := make(chan os.Signal, 1)
  signal.Notify(finish, syscall.SIGINT, syscall.SIGTERM)

  <-finish

  fmt.Println("\nClosing...")

  depthStream.Close()
  relay.Stop()
}

func main() {
  c := ParseConfig()

  if c.info {
    count := CountDevices()

    if(count == 1) {
      fmt.Printf("There is one Kinect connected.\n")
    } else {
      fmt.Printf("There are %d Kinects connected.\n", count)
    }
  } else {
    if c.device >= 0 && c.port > 100 {
      if c.device >= CountDevices() {
        fmt.Printf("Requested device %d is not connected!", c.device);
      } else {
        fmt.Printf("Start stream server with data from device %d...\n", c.device)
        start(c)
      }
    } else {
      fmt.Printf("Specify a device id >= 0 and port >= 100!\n")
    }
  }
}
