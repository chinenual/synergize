#!/bin/bash

iterations=${1:-20}

bash testfiles/bin/serialtests.sh ${iterations} output/linux-amd64/Synergize -port /dev/ttyS0 -baud 19200 2>&1 | tee serialtests-linux.log
