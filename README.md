# depthstream

**stream kinect's depth image to websocket clients**

With Depsthream you can easily get depth image data from your Kinect using a WebSocket Connection. A wepage running on the same computer can that way easily obtain the depth pixel array and feed into an animation.

## Usage

```
Depthstream.

Usage:
    depthstream info
    depthstream start <device> <port>
```

Make a connection to the server and request a frame or stream mode by sending a message:

```js
var ws = new WebSocket('ws://localhost:8080');
ws.binaryType = 'arraybuffer';

ws.onopen = function (event) {
  ws.send('1'); // request the current cashed frame
  ws.send('*'); // add connection to the stream list
};

ws.onmessage = function (message) {
  var array = new Uint8Array(message.data);
};
```
