package main

import (
    "code.google.com/p/go.net/websocket"
    "flag"
    "fmt" 
    "log"
    "net/http"
    "text/template"
)

////////////////////////////////////////
// Cht Connection
////////////////////////////////////////

type connection struct {
    // The websocket connection.
    ws *websocket.Conn

    // Buffered channel of outbound messages.
    send chan string
}

func (c *connection) reader() {
    for {
        var message string
        err := websocket.Message.Receive(c.ws, &message)
        if err != nil {
            break
        }
        // IF this is a chat message, broadcast it:
        h.broadcast <- message
        // TODO: IF this is a movement command, hand it to nav system
        // TODO: IF we don't recognize it, throw an error
    }
    c.ws.Close()
}

func (c *connection) writer() {
    for message := range c.send {
        fmt.Println("Sending message: " + message) 
        err := websocket.Message.Send(c.ws, message)
        if err != nil {
            break
        }
    }
    c.ws.Close()
}

func wsHandler(ws *websocket.Conn) {
    c := &connection{send: make(chan string, 256), ws: ws}
    h.register <- c
    defer func() { h.unregister <- c }()
    go c.writer()
    c.reader()
}

////////////////////////////////////////
// Chat Hub
////////////////////////////////////////

type hub struct {
    // Registered connections.
    connections map[*connection]bool

    // Inbound messages from the connections.
    broadcast chan string

    // Register requests from the connections.
    register chan *connection

    // Unregister requests from connections.
    unregister chan *connection
}

var h = hub{
    broadcast: make(chan string),
    register: make(chan *connection),
    unregister: make(chan *connection),
    connections: make(map[*connection]bool),
}

func (h *hub) run() {
    for {
        select {
        case c := <-h.register:
            fmt.Println("Client registered!")
            h.connections[c] = true
        case c := <-h.unregister:
            fmt.Println("Client unregistered!")
            delete(h.connections, c)
            close(c.send)
        case m := <-h.broadcast:
            fmt.Println("message broadcast!")
            for c := range h.connections {
                select {
                case c.send <- m:
                default:
                    delete(h.connections, c)
                    close(c.send)
                    go c.ws.Close()
                }
            }
        }
    }
}

////////////////////////////////////////
// Main
////////////////////////////////////////

var addr = flag.String("addr", ":8080", "http service address")
var homeTempl = template.Must(template.ParseFiles("home.html"))

func homeHandler(c http.ResponseWriter, req *http.Request) {
    homeTempl.Execute(c, req.Host)
}

func main() {
    flag.Parse()
    go h.run()
    http.HandleFunc("/", homeHandler)
    http.Handle("/ws", websocket.Handler(wsHandler))
    if err := http.ListenAndServe(*addr, nil); err != nil {
        log.Fatal("ListenAndServe:", err)
    }
}
