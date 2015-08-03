# depthstream

**stream kinect's depth image to websocket clients**

With Depsthream you can easily get depth image data from your Kinect using a WebSocket Connection. A wepage running on the same computer can that way easily obtain the depth pixel array with full 120 frames per second.

## Usage

```
Depthstream.

Usage:
    depthstream info
    depthstream start <device> <port>
```

Make a connection to the server and request a frame by sending an empty message.
