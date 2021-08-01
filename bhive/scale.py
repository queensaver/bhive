#! /usr/bin/python3

import time
import sys
import RPi.GPIO as GPIO
from hx711 import HX711
import argparse


hx = HX711(27, 22)
parser = argparse.ArgumentParser(description='beehive interface.')
parser.add_argument('--reference_unit', type=float, help='The reference unit we divide the measurement by to get the desired weight.', default=20.50671550671551)
parser.add_argument('--offset', type=float, help='The offset in grams we substract from the measurement to tare it.', default=8524115.625)
parser.add_argument('--calibrate', type=bool, help='Run the calibration.', default=False)
args = parser.parse_args()

hx.set_offset(args.offset)
hx.set_scale(args.reference_unit)

def setup():
    """
    code run once
    """
    print("Initializing.\n Please ensure that the scale is empty.")
    scale_ready = False
    while not scale_ready:
        if (GPIO.input(hx.DOUT) == 0):
            scale_ready = False
        if (GPIO.input(hx.DOUT) == 1):
            print("Initialization complete!")
            scale_ready = True

def calibrate():
    readyCheck = input("Remove any items from scale. Press any key when ready.")
    offset = hx.read_average()
    print("Value at zero (offset): {}".format(offset))
    hx.set_offset(offset)
    print("Please place an item of known weight on the scale.")

    readyCheck = input("Press any key to continue when ready.")
    measured_weight = (hx.read_average()-hx.get_offset())
    item_weight = input("Please enter the item's weight in grams.\n>")
    scale = int(measured_weight)/int(item_weight)
    hx.set_scale(scale)
    print("Scale adjusted for grams: {}".format(scale))

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
    if args.calibrate:
        setup()
        calibrate()
        GPIO.cleanup()
    try:
        weight = measureWeight()
        print(weight, end='')
    except OSError as e:
        sys.stderr.write("Failed to measure weight: %s" % e)
    finally:
        GPIO.cleanup()

main()
