#!/bin/bash

iterations=${1:-100}

bash testbin/serialtests.sh ${iterations} output/darwin-amd64/Synergize.app/Contents/MacOS/Synergize -port /dev/tty.usbserial-AL05OC8S -baud 9600 2>&1 | tee serialtests-mac.log
