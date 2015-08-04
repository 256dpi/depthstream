# depthstream

**stream kinect's depth image to websocket clients**

With Depsthream you can easily get depth image data from your Kinect using a WebSocket Connection. A wepage running on the same computer can that way easily obtain the depth pixel array and feed into an animation.

![Example](http://joel-github-static.s3.amazonaws.com/depthstream/screenshot.png)

## Usage

```
Depthstream.

Usage:
  depthstream [options]

Options:
  -h --help         Show this screen.
  -i --info         Show connected Kinects.
  -p --port=<num>   Port for server. [default: 9090].
  -d --dev=<id>     Device to open. [default: 0].
```

Make a connection to the server and request a frame or stream mode by sending a message:

```js
var ws = new WebSocket('ws://localhost:9090');
ws.binaryType = 'arraybuffer';

ws.onopen = function (event) {
  ws.send('1'); // request the current cashed frame
  ws.send('*'); // add connection to the stream list
};

ws.onmessage = function (message) {
  // the array contains 640*480 values from 0 to 255 representing
  // the detected depth from 0 to 10000mm
  var array = new Uint16Array(message.data);
};
```

See the [example](https://github.com/256dpi/depthstream/tree/master/test).

## Install (OSX only)

First install dependencies using brew:

```bash
$ brew install libusb
$ brew install libfreenect
```

Then download the latest binary from [releases](https://github.com/256dpi/depthstream/releases).

## Build

```bash
$ gopm get
$ gopm build
```
