#! /usr/bin/python3

import time
import sys
import RPi.GPIO as GPIO
from hx711 import HX711
import argparse
from uuid import getnode as get_mac
import requests


hx = HX711(5, 6)
mac = hex(get_mac())
parser = argparse.ArgumentParser(description='beehive interface.')
parser.add_argument('--reference_unit', type=float, help='The reference unit we divide the measurement by to get the desired weight.', default=20.544371)
parser.add_argument('--offset', type=int, help='The offset in grams we substract from the measurement to tare it.', default=2115)
parser.add_argument('--server_addr', type=str, help='HTTP server that implements the bBox REST API.', default="http://machine.intranet.wogri.com:8333/scale")
args = parser.parse_args()

hx.set_reading_format("MSB", "MSB")
hx.set_reference_unit(args.reference_unit)

def measureWeight():
    try:
        hx.reset()
        val = hx.get_weight(5) - args.offset
        if val < 0:
            val = 0
        print(val)
        hx.power_down()
    except Exception as e:
        sys.stderr.write(e)
    return val

def pushWeightToServer(weight):
    w = {"Weight": weight, "BBoxID": mac}
    resp = requests.post(args.server_addr, json=w)
    if resp.status_code != 200:
        sys.stderr.write("Something went terribly wrong, http status %s" % resp.status_code)

def run():
    try:
        weight = measureWeight()
        pushWeightToServer(weight)
    # if the server can't be reached - log.
    except OSError as e:
        sys.stderr.write("Failed to connect to %s: %s" %(args.server_addr, e))
    finally:
        GPIO.cleanup()

run()
