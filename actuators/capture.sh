#!/bin/bash
python3 /home/pi/hw.py --debug scan do light 1
sleep 1
# D=$(date +%s)
D=/home/pi/scan
libcamera-jpeg --brightness=0.1 -o $D.jpg
python3 /home/pi/hw.py --debug scan do light 0
curl -v -F bhiveId=b827eb3dcc6b -F epoch=$D -F scan=@$D.jpg http://localhost:80/varroa
python3 /home/pi/hw.py --debug scan do move 100 6
