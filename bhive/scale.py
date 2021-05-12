#! /usr/bin/python3

import time
import sys
import RPi.GPIO as GPIO
from hx711 import HX711
import argparse


hx = HX711(5, 6)
parser = argparse.ArgumentParser(description='beehive interface.')
parser.add_argument('--reference_unit', type=float, help='The reference unit we divide the measurement by to get the desired weight.', default=20.50671550671551)
parser.add_argument('--offset', type=float, help='The offset in grams we substract from the measurement to tare it.', default=8524115.625)
args = parser.parse_args()

hx.set_offset(args.offset)
hx.set_scale(args.reference_unit)

def measureWeight():
    try:
        val = hx.get_grams()
        if val < 0:
            val = 0
        hx.power_down()
    except Exception as e:
        sys.stderr.write(e)
    return val

def main():
    try:
        weight = measureWeight()
        print(weight, end='')
    except OSError as e:
        sys.stderr.write("Failed to measure weight: %s" % e)
    finally:
        GPIO.cleanup()

main()
