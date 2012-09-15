import socket
import threading

UDP_IP = "127.0.0.1"
UDP_PORT = 4004

BUFFER_SIZE = 256

HELLO = "Hello, World!"

print "UDP target port:", UDP_PORT

sock = socket.socket(socket.AF_INET, # Internet
                     socket.SOCK_DGRAM) # UDP

class ChatListenerThread(threading.Thread):
  socket = None
  #def __init__ (self, socket):
  #  self.socket = socket
  #  super(ChatListenerThread, self).__init__(self)

  def run (self):
    print "started listener"
    while True:
      string, address = self.socket.recvfrom(BUFFER_SIZE)
      print "Received message: %s" % message

class ChatWriterThread(threading.Thread):
  socket = None
  #def __init__(self, socket):
  #  self.socket = socket
  #  super(ChatWriterThread, self).__init__(self)

  def run (self):
    print "started writer"
    while True:
      message = raw_input('> ')
      #print "message:", message
      self.socket.sendto(message, (UDP_IP, UDP_PORT))

listener = ChatListenerThread()
writer   = ChatWriterThread()

listener.socket = sock
writer.socket = sock

listener.start()
writer.start()
