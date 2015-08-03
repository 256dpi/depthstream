package main

import (
  "github.com/gorilla/websocket"
  "time"
  "net/http"
  "fmt"
)

type connection struct {
  ws *websocket.Conn
  send chan []byte
  relay *Relay
}

func (c *connection) write(mt int, payload []byte) error {
  c.ws.SetWriteDeadline(time.Now().Add(time.Second))
  return c.ws.WriteMessage(mt, payload)
}

func (c *connection) writeLoop() {
  defer func(){
    c.ws.Close()
  }()

  for {
    select {
    case message, ok := <-c.send:
      if !ok {
        c.write(websocket.CloseMessage, []byte{})
        return
      }
      if err := c.write(websocket.BinaryMessage, message); err != nil {
        return
      }
    }
  }
}

func (c *connection) readLoop() {
  defer func() {
    c.relay.unregister <- c
    c.ws.Close()
  }()

  for {
    _, msg, err := c.ws.ReadMessage()
    if err != nil {
      break
    }

    str := string(msg)

    if str == "1" {
      c.relay.queue <- c
    } else if str == "*" {
      c.relay.stream <- c
    }
  }
}

type Relay struct {
  queue chan *connection
  stream chan *connection
  unstream chan *connection
  connections map[*connection]bool
  register chan *connection
  unregister chan *connection
  upgrader *websocket.Upgrader
}

func NewRelay() *Relay {
  return &Relay{
    queue: make(chan *connection),
    stream: make(chan *connection),
    unstream: make(chan *connection),
    connections: make(map[*connection]bool),
    register: make(chan *connection),
    unregister: make(chan *connection),
    upgrader: &websocket.Upgrader{
      ReadBufferSize: 1024,
      WriteBufferSize: 1024,
      CheckOrigin: func(r *http.Request) bool { return true },
    },
  }
}

func manage(r *Relay) {
  for {
    select {
    case c := <-r.register:
      r.connections[c] = true
      fmt.Printf("New client, total: %d\n", len(r.connections))
    case c := <-r.unregister:
      r.unstream <- c
      if _, ok := r.connections[c]; ok {
        delete(r.connections, c)
        close(c.send)
      }
      fmt.Printf("Lost client, total: %d\n", len(r.connections))
    }
  }
}

func (r *Relay) Start(port int) {
  go manage(r)

  http.HandleFunc("/", func(res http.ResponseWriter, req *http.Request){
    if req.Method != "GET" {
      http.Error(res, "Method not allowed", 405)
      return
    }

    ws, err := r.upgrader.Upgrade(res, req, nil)
    if err != nil {
      panic(err)
    }

    c := &connection{
      send: make(chan []byte, 256),
      ws: ws,
      relay: r,
    }

    r.register <- c

    go c.writeLoop()
    go c.readLoop()
  })

  go func(){
    err := http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", port), nil)
    if err != nil {
      panic(err)
    }
  }()

  fmt.Printf("Server launched on port %d!\n", port)
}

func (r *Relay) Stop() {
  //TODO: disconnect clients
}
