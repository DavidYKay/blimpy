#!/usr/bin/env python

import pygame

VALID_JOYSTICK_AXES = [
    0,
    1,
]

class App:
  x = y = 0
  def __init__(self):
    pygame.display.set_caption("Blimpy Client")

    self.init_joystick()

    #x = y = 0
    self.running = 1
    self.screen = pygame.display.set_mode((640, 400))
    #self.font = pygame.font.SysFont("Courier", 20)
  
  def init_joystick(self):
    self.my_joystick = None
    self.joystick_names = []
    
    pygame.joystick.init()

    for i in range(0, pygame.joystick.get_count()):
      self.joystick_names.append(pygame.joystick.Joystick(i).get_name())

    print self.joystick_names
    if (len(self.joystick_names) > 0):
      self.my_joystick = pygame.joystick.Joystick(0)
      self.my_joystick.init()

  def main(self):
    while self.running:
      event = pygame.event.poll()
      if event.type == pygame.QUIT:
        self.running = 0
      elif event.type == pygame.MOUSEMOTION:
        #print "mouse at (%d, %d)" % event.pos
        pass
      elif event.type == pygame.MOUSEBUTTONDOWN:
        delta = pygame.mouse.get_rel()
        print "mouse moved (%d, %d)" % delta
      elif event.type == pygame.MOUSEBUTTONUP:
        pass
      elif event.type == pygame.JOYAXISMOTION:
        if event.axis in VALID_JOYSTICK_AXES:
          print "joystick axis %d moved %f" % (event.axis, event.value,)
        else:
          pass

      self.screen.fill((0, 0, 0))
      pygame.display.flip()

app = App()
app.main()
