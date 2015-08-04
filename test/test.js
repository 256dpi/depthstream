var info = document.getElementById('info');
var canvas = document.getElementById('depth');

var ctx = canvas.getContext('2d');

var frames = 0;
var bytes = 0;

var width = 640 / 1;
var height = 480 / 1;

canvas.width = width;
canvas.height = height;

function process(array){
  var img = ctx.createImageData(width, height);

  for(var i=0; i<height; i++) {
    for(var j=0; j<width; j++) {
      var pos = i * width + j;
      var index = pos * 4;

      if(array[pos] == 0) {
        img.data[index] = 255;
        img.data[index + 1] = 0;
        img.data[index + 2] = 0;
        img.data[index + 3] = 255;
      } else {
        var value = 255 - (array[pos] / 10000 * 255);
        img.data[index] = value;
        img.data[index + 1] = value;
        img.data[index + 2] = value;
        img.data[index + 3] = 255;
      }
    }
  }

  ctx.putImageData(img, 0, 0);
};

var ws = new WebSocket('ws://localhost:9090');
ws.binaryType = 'arraybuffer';

ws.onopen = function (event) {
  ws.send('*');
};

ws.onmessage = function (message) {
  frames++;
  bytes += message.data.byteLength;
  process(new Uint16Array(message.data));
};

setInterval(function(){
  var mbs =  Math.round(bytes / 1024 / 1024 * 100) / 100;
  info.innerHTML = frames + ' F/s - ' + mbs + ' MB/s';
  frames = 0;
  bytes = 0;
}, 1000);
