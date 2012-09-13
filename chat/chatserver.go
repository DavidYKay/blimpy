package main

import (
    //"code.google.com/p/go.net/websocket"
    "bytes"
    "flag"
    "fmt" 
    "log"
    "net/http"
    "net"
    "text/template"
)

////////////////////////////////////////
// Chat Connection
////////////////////////////////////////

const BUFFER_SIZE = 256

//type wsconnection struct {
//    // The websocket connection.
//    socket *websocket.Conn
//
//    // Buffered channel of outbound messages.
//    send chan string
//}

//type connection struct {
//    // The websocket connection.
//    socket *net.Conn
//
//    // Buffered channel of outbound messages.
//    send chan string
//}

//type connection struct {
type connection struct {
    // The websocket connection.
    socket *net.UDPConn

    // Buffered channel of outbound messages.
    send chan string
}

//func (c *connection) wsReader() {
//    for {
//        var message string
//        err := websocket.Message.Receive(c.socket, &message)
//        if err != nil {
//            break
//        }
//        // IF this is a chat message, broadcast it:
//        h.broadcast <- message
//        // TODO: IF this is a movement command, hand it to nav system
//        // TODO: IF we don't recognize it, throw an error
//    }
//    c.socket.Close()
//}

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


//func (c *connection) wsWriter() {
//    for message := range c.send {
//        fmt.Println("Sending message: " + message) 
//        err := websocket.Message.Send(c.socket, message)
//        if err != nil {
//            break
//        }
//    }
//    c.socket.Close()
//}

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


//func handleRawSocket(socket net.Conn) {
//    c := &connection{send: make(chan string, 256), socket: socket}
//    h.register <- c
//    defer func() { h.unregister <- c }()
//    go c.writer()
//    c.reader()
//}

//func wsHandler(ws *websocket.Conn) {
//    c := &connection{send: make(chan string, BUFFER_SIZE), socket: ws}
//    h.register <- c
//    defer func() { h.unregister <- c }()
//    go c.writer()
//    c.wsReader()
//}

//func handleTcp(socket TCPConn) {
//    handleRawSocket(socket)
//}

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

var addr = flag.String("addr", ":8080", "http service address")
var homeTempl = template.Must(template.ParseFiles("home.html"))

func homeHandler(c http.ResponseWriter, req *http.Request) {
    homeTempl.Execute(c, req.Host)
}

func main() {
    flag.Parse()
    go h.run()

    //http.HandleFunc("/", homeHandler)
    //http.Handle("/ws", websocket.Handler(wsHandler))
    //if err := http.ListenAndServe(*addr, nil); err != nil {
    //    log.Fatal("Websocket ListenAndServe:", err)
    //}

    //ln, err := net.Listen("tcp", ":5005")
    //if err != nil {
    //    log.Fatal("TCP Listen:", err)
    //}
    //for {
    //    conn, err := ln.Accept()
    //    if err != nil {
    //        log.Error("TCP accept error:", err)
    //        // handle error
    //        continue
    //    }
    //    go handleTcp(conn)
    //}

    //ln, err := net.ListenUDP("udp", ":4004")
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
