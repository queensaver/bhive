#! /usr/bin/python3

import time
import sys
import RPi.GPIO as GPIO
from hx711 import HX711
from http.server import BaseHTTPRequestHandler, HTTPServer
import argparse

hx = HX711(5, 6)
parser = argparse.ArgumentParser(description='beehive interface.')
parser.add_argument('--reference_unit', type=float, help='The reference unit we divide the measurement by to get the desired weight.', default=20.544371)
parser.add_argument('--offset', type=int, help='The offset in grams we substract from the measurement to tare it.', default=2115)
parser.add_argument('--http_port', type=int, help='The http port to listen on (default: 8000)', default=8000)
args = parser.parse_args()

def cleanAndExit():
    GPIO.cleanup()
    sys.exit()

hx.set_reading_format("MSB", "MSB")
hx.set_reference_unit(args.reference_unit)

class RequestHandler(BaseHTTPRequestHandler):
    def do_GET(self):
        msg = "unkown error"
        return_code = 500
        try:
            hx.reset()
            val = hx.get_weight(5) - args.offset
            if val < 0:
                val = 0
            print(val)
            msg = str(val)
            return_code = 200

            hx.power_down()
            hx.power_up()
        except Exception as e:
            sys.stderr.write(e)

        self.protocol_version = "HTTP/1.1"
        self.send_response(return_code)
        self.send_header("Content-Length", len(msg))
        self.end_headers()
        self.wfile.write(bytes(msg, "utf8"))

def run():
    server = ('', args.http_port)
    httpd = HTTPServer(server, RequestHandler)
    httpd.serve_forever()

run()
