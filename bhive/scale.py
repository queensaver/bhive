#! /usr/bin/python3

import time
import sys
import RPi.GPIO as GPIO
from hx711 import HX711
import argparse


hx = HX711(5, 6)
parser = argparse.ArgumentParser(description='beehive interface.')
parser.add_argument('--reference_unit', type=float, help='The reference unit we divide the measurement by to get the desired weight.', default=20.544371)
parser.add_argument('--offset', type=int, help='The offset in grams we substract from the measurement to tare it.', default=2115)
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

def run():
    try:
        weight = measureWeight()
        print(weight)
    except OSError as e:
        sys.stderr.write("Failed to measure weight: %s" % e)
    finally:
        GPIO.cleanup()

run()
