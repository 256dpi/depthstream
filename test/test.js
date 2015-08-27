var info = document.getElementById('info');
var depth = document.getElementById('depth');
var color = document.getElementById('color');

var depthContext = depth.getContext('2d');
var colorContext = color.getContext('2d');

var frames = 0;
var bytes = 0;

var width = 640 / 1;
var height = 480 / 1;

depth.width = width;
depth.height = height;

color.width = width;
color.height = height;

function process(array16, array8){
  var depthImage = depthContext.createImageData(width, height);
  var colorImage = colorContext.createImageData(width, height);

  for(var i=0; i<height; i++) {
    for(var j=0; j<width; j++) {
      var pos = i * width + j;
      var pixel = pos * 4;

      if(array16[pos] == 0) {
        depthImage.data[pixel] = 255;
        depthImage.data[pixel + 1] = 0;
        depthImage.data[pixel + 2] = 0;
        depthImage.data[pixel + 3] = 255;
      } else {
        var value = 255 - (array16[pos] / 10000 * 255);
        depthImage.data[pixel] = value;
        depthImage.data[pixel + 1] = value;
        depthImage.data[pixel + 2] = value;
        depthImage.data[pixel + 3] = 255;
      }
    }
  }
  
  var offset = width * height * 2;

  if(array8.length > offset) {
    for(var i=0; i<height; i++) {
      for(var j=0; j<width; j++) {
        var pos = i * width + j;
        var pixel = pos * 4;
        var source = pos * 3;

        colorImage.data[pixel] = array8[offset + source];
        colorImage.data[pixel + 1] = array8[offset + source + 1];
        colorImage.data[pixel + 2] = array8[offset + source + 2];
        colorImage.data[pixel + 3] = 255;
      }
    }
  }

  depthContext.putImageData(depthImage, 0, 0);
  colorContext.putImageData(colorImage, 0, 0);
};

var ws = new WebSocket('ws://localhost:9090');
ws.binaryType = 'arraybuffer';

ws.onopen = function (event) {
  ws.send('*');
};

ws.onmessage = function (message) {
  frames++;
  bytes += message.data.byteLength;
  process(new Uint16Array(message.data), new Uint8Array(message.data));
};

setInterval(function(){
  var mbs =  Math.round(bytes / 1024 / 1024 * 100) / 100;
  info.innerHTML = frames + ' F/s - ' + mbs + ' MB/s';
  frames = 0;
  bytes = 0;
}, 1000);
