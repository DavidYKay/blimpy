package main

import (
    //"code.google.com/p/go.net/websocket"
    "bytes"
    "flag"
    "fmt" 
    "log"
    "net"
)

////////////////////////////////////////
// Chat Connection
////////////////////////////////////////

const BUFFER_SIZE = 256

//type connection struct {
type connection struct {
    // The websocket connection.
    socket *net.UDPConn

    // Buffered channel of outbound messages.
    send chan string
}

func (c *connection) udpReader() {
    for {
        var message string

        //func (c *UDPConn) ReadFromUDP(b []byte) (n int, addr *UDPAddr, err error)
        
        buf := make([]byte, 1024)
        n, err := c.socket.Read(buf)

        if err != nil {
            log.Fatalln("error reading UDP: ", err)
            //break
	}
        message = string(buf[0:n])

        // IF this is a chat message, broadcast it:
        h.broadcast <- message
        // TODO: IF this is a movement command, hand it to nav system
        // TODO: IF we don't recognize it, throw an error
    }
    c.socket.Close()
}

func (c *connection) udpWriter() {
    for message := range c.send {
        fmt.Println("Sending message: " + message) 

        // TODO: convert to bytes
        //buf := strings.Buffer.NewBufferString(message)
        buf := bytes.NewBufferString(message).Bytes()

        n, err := c.socket.Write(buf) 

        if err != nil {
            log.Fatalln("error writing to UDP: ", err, n)
            //log.Fatalln(n)
            //log.Fatalln(err)
            //break
	}

    }
    c.socket.Close()
}

func handleUdp(socket *net.UDPConn) {
    c := &connection{send: make(chan string, BUFFER_SIZE), socket: socket}
    h.register <- c
    defer func() { h.unregister <- c }()
    go c.udpWriter()
    c.udpReader()
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
                    go c.socket.Close()
                }
            }
        }
    }
}

////////////////////////////////////////
// Main
////////////////////////////////////////

func main() {
    flag.Parse()
    go h.run()

    udpAddr := &net.UDPAddr{ IP: net.ParseIP("127.0.0.1"), Port: 4004 }
    conn, err := net.ListenUDP("udp", udpAddr)
    if err != nil {
        log.Fatal("UDP Listen:", err)
    } else {
        fmt.Println("listening on ", conn.LocalAddr().String())
    }

    handleUdp(conn)

    //for {
    //    conn, err := ln.Accept()
    //    if err != nil {
    //        log.Fatal("UDP accept error:", err)
    //        // handle error
    //        continue
    //    }
    //    go handleUdp(conn)
    //}

}
