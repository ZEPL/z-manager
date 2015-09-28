#!/usr/bin/env python
# -*- coding: utf-8 -*-
# Finds an open port and prints it to stdio

import socket
import SocketServer
import SimpleHTTPServer

def find():
  s = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
  s.bind(('', 0))
  addr = s.getsockname()
  s.close()
  return addr[1]


class Handler(SimpleHTTPServer.SimpleHTTPRequestHandler):
  def do_GET(self):
    self.send_response(200)
    self.end_headers()
    self.wfile.write(find())

def main():
  httpd = SocketServer.ForkingTCPServer(('', 7777), Handler)
  httpd.serve_forever()

if __name__ == '__main__':
  main()
