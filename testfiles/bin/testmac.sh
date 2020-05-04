#!/bin/bash

iterations=${1:-20}

bash testfiles/bin/serialtests.sh ${iterations} output/darwin-amd64/Synergize.app/Contents/MacOS/Synergize -port /dev/tty.usbserial-AL05OC8S 2>&1 | tee serialtests-mac.log
