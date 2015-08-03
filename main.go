package main

import (
  "fmt"
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
      StreamDepthImage(c.port, c.device)
    } else {
      fmt.Printf("Specify a device id >= 0 and port >= 100!\n")
    }
  }
}
