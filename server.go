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

func (c *connection) loop() {
  defer func(){
    c.ws.Close()
    c.relay.unregister <- c
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

type Relay struct {
  connections map[*connection]bool
  broadcast chan []byte
  register chan *connection
  unregister chan *connection
  upgrader *websocket.Upgrader
}

func NewRelay() *Relay {
  return &Relay{
    connections: make(map[*connection]bool),
    broadcast: make(chan []byte),
    register: make(chan *connection),
    unregister: make(chan *connection),
    upgrader: &websocket.Upgrader{
      ReadBufferSize: 1024,
      WriteBufferSize: 1024,
      CheckOrigin: func(r *http.Request) bool { return true },
    },
  }
}

func relay(r *Relay) {
  for {
    select {
    case c := <-r.register:
      r.connections[c] = true
      fmt.Printf("New client, total: %d\n", len(r.connections))
    case c := <-r.unregister:
      if _, ok := r.connections[c]; ok {
        delete(r.connections, c)
        close(c.send)
      }
      fmt.Printf("Lost client, total: %d\n", len(r.connections))
    case m := <-r.broadcast:
      for c := range r.connections {
        c.send <- m
      }
    }
  }
}

func (r *Relay) Start(port int) {
  go relay(r)

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

    go c.loop()
  })

  go func(){
    err := http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", port), nil)
    if err != nil {
      panic(err)
    }
  }()

  fmt.Printf("Server launched on port %d\n", port)
}

func (r *Relay) Forward(msg []byte) {
  r.broadcast <- msg
}

func (r *Relay) Stop() {
  // does nothing at the moment
}
