package main

/*
 * ASYNCHRONOUS
 * Event-based 
 */
func handleMessage(Message message) {
    if (message.type == Command.MOVE) {

    } else if (message.type == Command.SHOOT) {

    }
}

func main() {

    Socket socket = new Socket('localhost', 1337);
    socket.attachListener(handleMessage);
    socket.listen();
    
}
