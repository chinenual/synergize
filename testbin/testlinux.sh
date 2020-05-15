#!/bin/bash

iterations=${1:-100}

bash testbin/serialtests.sh ${iterations} output/linux-amd64/Synergize -port /dev/ttyS0 -baud 9600 2>&1 | tee serialtests-linux.log
